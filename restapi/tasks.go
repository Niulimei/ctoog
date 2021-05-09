package restapi

import (
	"bytes"
	"ctgb/database"
	"ctgb/models"
	"ctgb/restapi/operations"
	"ctgb/utils"
	"encoding/json"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func init() {
	go func() {
		t := time.NewTicker(time.Second * 5)
		for {
			select {
			case <-t.C:
				log.Debug("ticker begin")
				var taskLogs []*database.TaskLog
				now := time.Now()
				start := now.Add(time.Hour * -1).Format("2006-01-02 15:04:05")
				log.Debug("start", start)
				log.Debug(database.DB.Select(&taskLogs,
					"SELECT * FROM task_log WHERE start_time < $1 AND status = 'running'", start))
				tx, err := database.DB.Begin()
				if err == nil {
					for _, taskLog := range taskLogs {
						log.Debug("auto close task log ", taskLog.LogId)
						tx.Exec("UPDATE task_log SET status = 'failed', end_time = $1 WHERE log_id = $2",
							now.Format("2006-01-02 15:04:05"), taskLog.LogId)
						tx.Exec("UPDATE task SET status = 'failed' WHERE id = $1", taskLog.TaskId)
					}
					tx.Commit()
				} else {
					log.Error(err)
				}
			}
		}
	}()
}

func startTask(taskID int64) {
	task := &database.TaskModel{}
	err := database.DB.Get(task, "SELECT cc_password,"+
		" cc_user, component, git_password, git_url, git_user, git_email, pvob, include_empty"+
		" FROM task WHERE id = $1", taskID)
	if err != nil {
		log.Error("start task but db err:", err)
		return
	}
	worker := &database.WorkerModel{}
	if task.WorkerId != 0 {
		err = database.DB.Get(worker, "SELECT * FROM worker WHERE id = $1", task.WorkerId)
	} else {
		err = database.DB.Get(worker, "SELECT * FROM worker ORDER BY task_count DESC limit 1")
	}
	workerUrl := worker.WorkerUrl
	if worker.WorkerUrl == "" {
		log.Error("get worker with no url:", worker)
		return
	}
	var matchInfo []*models.TaskMatchInfo
	database.DB.Select(&matchInfo, "SELECT git_branch, stream FROM match_info WHERE task_id = $1 ORDER BY id",
		taskID)
	startTime := time.Now().Format("2006-01-02 15:04:05")
	r := database.DB.MustExec(
		"INSERT INTO task_log (task_id, status, start_time, end_time, duration)"+
			" VALUES($1, 'running', $2, $3, 0)", taskID, startTime, "",
	)
	taskLogId, err := r.LastInsertId()
	if err == nil {
		type InnerMatchInfo struct {
			Branch string
			Stream string
		}

		type InnerTask struct {
			TaskId       int64
			TaskLogId    int64
			CcPassword   string
			CcUser       string
			Component    string
			GitPassword  string
			GitURL       string
			GitUser      string
			GitEmail     string
			Pvob         string
			IncludeEmpty bool
			Matches      []InnerMatchInfo
		}
		workerTaskModel := InnerTask{
			TaskId:       taskID,
			TaskLogId:    taskLogId,
			CcPassword:   task.CcPassword,
			CcUser:       task.CcUser,
			Component:    task.Component,
			GitPassword:  task.GitPassword,
			GitURL:       task.GitURL,
			GitUser:      task.GitUser,
			GitEmail:     task.GitEmail,
			Pvob:         task.Pvob,
			IncludeEmpty: task.IncludeEmpty,
		}
		for _, match := range matchInfo {
			workerTaskModel.Matches =
				append(workerTaskModel.Matches, InnerMatchInfo{Stream: match.Stream, Branch: match.GitBranch})
		}
		workerTaskModelByte, _ := json.Marshal(workerTaskModel)
		req, _ := http.NewRequest(http.MethodPost, "http://"+workerUrl+"/new_task", bytes.NewBuffer(workerTaskModelByte))
		req.Header.Set("Content-Type", "application/json")
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != http.StatusCreated {
			log.Error(fmt.Errorf("不能发送任务给%d", worker.Id), err)
			database.DB.MustExec("UPDATE task_log SET status = 'failed' WHERE id = $1", taskLogId)
			return
		}
	} else {
		log.Error(err)
		return
	}

	tx := database.DB.MustBegin()
	tx.MustExec(
		"UPDATE task SET status = 'running', worker_id = $1 WHERE id = $2", worker.Id, taskID,
	)
	tx.MustExec(
		"UPDATE worker SET task_count = task_count + 1 WHERE id = $1", worker.Id,
	)
	tx.Commit()
	return
}

func CreateTaskHandler(params operations.CreateTaskParams) middleware.Responder {
	userToken := params.Authorization
	username, verified := utils.Verify(userToken)
	if !verified {
		return middleware.Error(401, "鉴权失败")
	}
	taskInfo := params.TaskInfo
	r := database.DB.MustExec("INSERT INTO task (pvob, component, cc_user, cc_password, git_url,"+
		"git_user, git_password, status, last_completed_date_time, creator, include_empty, git_email)"+
		" VALUES ($1, $2, $3, $4, $5, $6, $7, $8, '', $9, $10, $11)",
		taskInfo.Pvob, taskInfo.Component, taskInfo.CcUser, taskInfo.CcPassword, taskInfo.GitURL,
		taskInfo.GitUser, taskInfo.GitPassword, "init", username,
		taskInfo.IncludeEmpty, taskInfo.GitEmail)
	taskId, err := r.LastInsertId()
	if err != nil {
		return operations.NewCreateTaskInternalServerError().WithPayload(
			&models.ErrorModel{Message: fmt.Sprintf("Insert into db error: %+v", err), Code: 500})
	}
	tx, _ := database.DB.Begin()
	for _, match := range taskInfo.MatchInfo {
		tx.Exec("INSERT INTO "+
			"match_info (task_id, stream, git_branch) "+
			"VALUES($1, $2, $3)",
			taskId, match.Stream, match.GitBranch)
	}
	tx.Commit()
	go startTask(taskId)
	return operations.NewCreateTaskCreated().WithPayload(&models.OK{Message: "ok"})
}

func GetTaskHandler(params operations.GetTaskParams) middleware.Responder {
	taskID := params.ID
	task := &models.TaskModel{}
	database.DB.Get(task, "SELECT cc_password,"+
		" cc_user, component, git_password, git_url, git_user, pvob, include_empty, git_email"+
		" FROM task WHERE id = $1", taskID)
	var matchInfo []*models.TaskMatchInfo
	database.DB.Select(&matchInfo, "SELECT git_branch, stream FROM match_info WHERE task_id = $1", taskID)
	task.MatchInfo = matchInfo
	var logList []*models.TaskLogInfo
	database.DB.Select(&logList, "SELECT duration, end_time, log_id, start_time, status FROM task_log WHERE task_id = $1 ORDER BY log_id", taskID)
	taskDetail := &models.TaskDetail{TaskModel: task, LogList: logList}
	return operations.NewGetTaskOK().WithPayload(taskDetail)
}

func ListTaskHandler(params operations.ListTaskParams) middleware.Responder {
	username, verified := utils.Verify(params.Authorization)
	if !verified {
		return middleware.Error(http.StatusUnauthorized, "鉴权失败")
	}
	var query, queryCount string
	user := getUserInfo(username)
	if user.RoleID == int64(AdminRole) {
		query = "SELECT pvob, component, git_url, id, last_completed_date_time," +
			" status, include_empty, git_email" +
			" FROM task WHERE creator = $1 or 1 = 1 ORDER BY id LIMIT $2 OFFSET $3;"
		queryCount = "SELECT count(id) FROM task;"
	} else {
		query = "SELECT pvob, component, git_url, id, last_completed_date_time," +
			" status, include_empty, git_email" +
			" FROM task WHERE creator = $1 ORDER BY id LIMIT $2 OFFSET $3;"
		queryCount = "SELECT count(id) FROM task WHERE creator = $1;"
	}
	var tasks []*models.TaskInfoModel
	var count int64
	err := database.DB.Select(&tasks, query, username, params.Limit, params.Offset)
	if err != nil {
		return middleware.Error(http.StatusInternalServerError, "Sql Error")
	}
	err = database.DB.Get(&count, queryCount, username)
	if err != nil {
		return middleware.Error(http.StatusInternalServerError, "Sql Error")
	}
	tasksPage := &models.TaskPageInfoModel{}
	tasksPage.TaskInfo = tasks
	tasksPage.Count = count
	return operations.NewListTaskOK().WithPayload(tasksPage)
}

func UpdateTaskHandler(params operations.UpdateTaskParams) middleware.Responder {
	//username, verified := utils.Verify(params.Authorization)
	taskId := params.ID
	task := &database.TaskModel{}
	database.DB.Get(task, "SELECT status, worker_id FROM task WHERE id = $1", taskId)
	taskLogInfo := params.TaskLog
	taskLog := &database.TaskLog{}
	err := database.DB.Get(taskLog, "SELECT * FROM task_log WHERE task_id = $1 AND status = 'running'", taskId)
	if err != nil {
		log.Error(err)
		return middleware.Error(404, "没发现任务")
	}
	tx := database.DB.MustBegin()
	log.Debug("update task:", params.TaskLog)
	if params.TaskLog.LogID != "" {
		tx.MustExec("UPDATE task_log SET status = $1, end_time = $2, duration = $3 WHERE log_id = $4",
			taskLogInfo.Status, taskLogInfo.EndTime, taskLogInfo.Duration, params.TaskLog.LogID)
		tx.MustExec("UPDATE task SET status = $1, last_completed_date_time = $2 WHERE id = $3",
			taskLogInfo.Status, taskLogInfo.EndTime, taskId)
		tx.MustExec("UPDATE worker SET task_count = task_count - 1 WHERE id = $1", task.WorkerId)
	} else {
		if params.TaskLog.Pvob != "" {
			tx.MustExec("UPDATE task SET pvob = $1 WHERE id = $2", params.TaskLog.Pvob, taskId)
		}
		if params.TaskLog.Component != "" {
			tx.MustExec("UPDATE task SET component = $1 WHERE id = $2", params.TaskLog.Component, taskId)
		}
		if params.TaskLog.IncludeEmpty != task.IncludeEmpty {
			tx.MustExec("UPDATE task SET include_empty = $1 WHERE id = $2", params.TaskLog.IncludeEmpty, taskId)
		}
		if len(params.TaskLog.MatchInfo) > 0 {
			tx.MustExec("DELETE FROM match_info WHERE task_id = $1", taskId)
			for _, match := range params.TaskLog.MatchInfo {
				tx.MustExec("INSERT INTO "+
					"match_info (task_id, stream, git_branch) "+
					"VALUES($1, $2, $3)",
					taskId, match.Stream, match.GitBranch)
			}
		}
	}
	log.Debug("task update commit:", tx.Commit())
	return operations.NewUpdateTaskCreated().WithPayload(&models.OK{
		Message: "ok",
	})
}
func RestartTaskHandler(params operations.RestartTaskParams) middleware.Responder {
	//username, verified := utils.Verify(params.Authorization)
	taskId := params.RestartTrigger.ID
	task := &database.TaskModel{}
	database.DB.Get(task, "SELECT status, worker_id FROM task WHERE id = $1", taskId)
	if task.Status != "running" {
		go startTask(taskId)
	}
	return operations.NewUpdateTaskCreated().WithPayload(&models.OK{
		Message: "ok",
	})
}
