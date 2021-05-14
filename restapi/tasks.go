package restapi

import (
	"bytes"
	"ctgb/database"
	"ctgb/models"
	"ctgb/restapi/operations"
	"ctgb/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-openapi/runtime/middleware"
	log "github.com/sirupsen/logrus"
)

func init() {
	go func() {
		t := time.NewTicker(time.Second * 5)
		for {
			select {
			case <-t.C:
				log.Info("ticker begin")
				var taskLogs []*database.TaskLog
				now := time.Now()
				start := now.Add(time.Hour * -1).Format("2006-01-02 15:04:05")
				log.Info("start", start)
				log.Info(database.DB.Select(&taskLogs,
					"SELECT * FROM task_log WHERE start_time < $1 AND status = 'running'", start))
				tx, err := database.DB.Begin()
				if err == nil {
					for _, taskLog := range taskLogs {
						log.Info("auto close task log ", taskLog.LogId)
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

func startTask(taskId int64) {
	task := &database.TaskModel{}
	err := database.DB.Get(task, "SELECT cc_password,"+
		" cc_user, component, git_password, git_url, git_user, git_email, pvob, include_empty, dir, keep, worker_id"+
		" FROM task WHERE id = $1", taskId)
	if err != nil {
		log.Error("start task but db err:", err)
		return
	}
	worker := &database.WorkerModel{}
	newAssigned := false
	log.Info(task.WorkerId)
	if task.WorkerId != 0 {
		log.Debug("got worker id")
		err = database.DB.Get(worker, "SELECT * FROM worker WHERE id = $1 and status = 'running'", task.WorkerId)
	} else {
		newAssigned = true
		err = database.DB.Get(worker, "SELECT * FROM worker WHERE status = 'running' ORDER BY task_count limit 1")
	}
	workerUrl := worker.WorkerUrl
	if worker.WorkerUrl == "" {
		log.Error("get worker with no url:", worker, err)
		return
	}
	var matchInfo []*models.TaskMatchInfo
	database.DB.Select(&matchInfo, "SELECT git_branch, stream FROM match_info WHERE task_id = $1 ORDER BY id",
		taskId)
	startTime := time.Now().Format("2006-01-02 15:04:05")
	r := database.DB.MustExec(
		"INSERT INTO task_log (task_id, status, start_time, end_time, duration)"+
			" VALUES($1, 'running', $2, $3, 0)", taskId, startTime, "",
	)
	taskLogId, err := r.LastInsertId()
	component := task.Component + task.Dir
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
			Keep         string
		}
		workerTaskModel := InnerTask{
			TaskId:       taskId,
			TaskLogId:    taskLogId,
			CcPassword:   task.CcPassword,
			CcUser:       task.CcUser,
			Component:    component,
			GitPassword:  task.GitPassword,
			GitURL:       task.GitURL,
			GitUser:      task.GitUser,
			GitEmail:     task.GitEmail,
			Pvob:         task.Pvob,
			IncludeEmpty: task.IncludeEmpty,
			Keep:         task.Keep,
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
		utils.RecordLog(utils.Info, utils.StartTask, "", fmt.Sprintf("TaskId: %d, TaskLogId: %d", taskId, taskLogId), 0)
		if err != nil || resp.StatusCode != http.StatusCreated {
			log.Error(fmt.Errorf("不能发送任务给%d", worker.Id), err)
			database.DB.MustExec("UPDATE task_log SET status = 'failed' WHERE log_id = $1", taskLogId)
			database.DB.MustExec("UPDATE task SET status = 'failed' WHERE id = $1", taskId)
			database.DB.MustExec("UPDATE worker SET status = 'dead' WHERE id = $1", worker.Id)
			return
		}
	} else {
		log.Error(err)
		return
	}

	if newAssigned {
		tx := database.DB.MustBegin()
		tx.MustExec(
			"UPDATE task SET status = 'running', worker_id = $1 WHERE id = $2", worker.Id, taskId,
		)
		tx.MustExec(
			"UPDATE worker SET task_count = task_count + 1 WHERE id = $1", worker.Id,
		)
		tx.Commit()
	}
	return
}

func CreateTaskHandler(params operations.CreateTaskParams) middleware.Responder {
	userToken := params.Authorization
	username, verified := utils.Verify(userToken)
	if !verified {
		return middleware.Error(401, models.ErrorModel{Message: "鉴权失败"})
	}
	taskInfo := params.TaskInfo
	if len(taskInfo.Dir) > 0 && !strings.HasPrefix(taskInfo.Dir, "/") {
		taskInfo.Dir = "/" + taskInfo.Dir
	}
	r := database.DB.MustExec("INSERT INTO task (pvob, component, cc_user, cc_password, git_url,"+
		"git_user, git_password, status, last_completed_date_time, creator, include_empty, git_email, dir, keep, worker_id)"+
		" VALUES ($1, $2, $3, $4, $5, $6, $7, $8, '', $9, $10, $11, $12, $13, 0)",
		taskInfo.Pvob, taskInfo.Component, taskInfo.CcUser, taskInfo.CcPassword, taskInfo.GitURL,
		taskInfo.GitUser, taskInfo.GitPassword, "init", username,
		taskInfo.IncludeEmpty, taskInfo.GitEmail, taskInfo.Dir, taskInfo.Keep)
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
	utils.RecordLog(utils.Info, utils.AddTask, "", fmt.Sprintf("TaskId: %d", taskId), 0)
	return operations.NewCreateTaskCreated().WithPayload(&models.OK{Message: "ok"})
}

func GetTaskHandler(params operations.GetTaskParams) middleware.Responder {
	taskID := params.ID
	task := &models.TaskModel{}
	database.DB.Get(task, "SELECT cc_password,"+
		" cc_user, component, git_password, git_url, git_user, pvob, include_empty, git_email, dir, keep"+
		" FROM task WHERE id = $1", taskID)
	var matchInfo []*models.TaskMatchInfo
	database.DB.Select(&matchInfo, "SELECT git_branch, stream FROM match_info WHERE task_id = $1", taskID)
	task.MatchInfo = matchInfo
	var logList []*models.TaskLogInfo
	database.DB.Select(&logList, "SELECT duration, end_time, log_id, start_time, status FROM task_log WHERE task_id = $1 ORDER BY log_id DESC", taskID)
	taskDetail := &models.TaskDetail{TaskModel: task, LogList: logList}
	return operations.NewGetTaskOK().WithPayload(taskDetail)
}

func ListTaskHandler(params operations.ListTaskParams) middleware.Responder {
	username, verified := utils.Verify(params.Authorization)
	if !verified {
		return middleware.Error(http.StatusUnauthorized, models.ErrorModel{Message: "鉴权失败"})
	}
	var query, queryCount string
	user := getUserInfo(username)
	if user.RoleID == int64(AdminRole) {
		query = "SELECT pvob, component, git_url, id, last_completed_date_time," +
			" status, include_empty, git_email, dir, keep" +
			" FROM task WHERE creator = $1 or 1 = 1 ORDER BY id LIMIT $2 OFFSET $3;"
		queryCount = "SELECT count(id) FROM task;"
	} else {
		query = "SELECT pvob, component, git_url, id, last_completed_date_time," +
			" status, include_empty, git_email, dir, keep" +
			" FROM task WHERE creator = $1 ORDER BY id LIMIT $2 OFFSET $3;"
		queryCount = "SELECT count(id) FROM task WHERE creator = $1;"
	}
	var tasks []*models.TaskInfoModel
	var count int64
	err := database.DB.Select(&tasks, query, username, params.Limit, params.Offset)
	if err != nil {
		log.Error(err)
		return middleware.Error(http.StatusInternalServerError, models.ErrorModel{Message: "Sql Error"})
	}
	err = database.DB.Get(&count, queryCount, username)
	if err != nil {
		log.Error(err)
		return middleware.Error(http.StatusInternalServerError, models.ErrorModel{Message: "Sql Error"})
	}
	tasksPage := &models.TaskPageInfoModel{}
	tasksPage.TaskInfo = tasks
	tasksPage.Count = count
	return operations.NewListTaskOK().WithPayload(tasksPage)
}

func UpdateTaskHandler(params operations.UpdateTaskParams) middleware.Responder {
	//username, verified := utils.Verify(params.Authorization)
	taskId := params.ID
	log.Debug(taskId)
	tx := database.DB.MustBegin()
	log.Debug("update task:", params.TaskLog)
	if params.TaskLog.LogID != "" {
		task := &database.TaskModel{}
		err := database.DB.Get(task, "SELECT status, worker_id FROM task WHERE id = $1", taskId)
		taskLogInfo := params.TaskLog
		if err != nil {
			log.Error(err)
			return middleware.Error(404, models.ErrorModel{Message: "没发现任务"})
		}
		tx.MustExec("UPDATE task_log SET status = $1, end_time = $2, duration = $3 WHERE log_id = $4",
			taskLogInfo.Status, taskLogInfo.EndTime, taskLogInfo.Duration, params.TaskLog.LogID)
		tx.MustExec("UPDATE task SET status = $1, last_completed_date_time = $2 WHERE id = $3",
			taskLogInfo.Status, taskLogInfo.EndTime, taskId)
		//tx.MustExec("UPDATE worker SET task_count = task_count - 1 WHERE id = $1", task.WorkerId)
		utils.RecordLog(utils.Info, utils.UpdateTask, "", fmt.Sprintf("TaskId: %s", taskId), 0)
	} else {
		log.Debug("update params:", params.TaskLog)
		tx.MustExec("UPDATE task SET pvob = $1, component = $2, dir = $3, cc_user = $4, cc_password = $5, "+
			"git_url = $6, git_user = $7, git_password = $8, git_email = $9, include_empty = $10 WHERE id = $11",
			params.TaskLog.Pvob, params.TaskLog.Component, params.TaskLog.Dir, params.TaskLog.CcUser,
			params.TaskLog.CcPassword, params.TaskLog.GitURL, params.TaskLog.GitUser, params.TaskLog.GitPassword,
			params.TaskLog.GitEmail, params.TaskLog.IncludeEmpty, params.ID)
		if len(params.TaskLog.MatchInfo) > 0 {
			tx.MustExec("DELETE FROM match_info WHERE task_id = $1", taskId)
			for _, match := range params.TaskLog.MatchInfo {
				tx.MustExec("INSERT INTO "+
					"match_info (task_id, stream, git_branch) "+
					"VALUES($1, $2, $3)",
					taskId, match.Stream, match.GitBranch)
			}
		}
		utils.RecordLog(utils.Info, utils.UpdateTask, "", fmt.Sprintf("TaskId: %s", taskId), 0)
		taskIdInt, _ := strconv.ParseInt(taskId, 10, 64)
		go startTask(taskIdInt)
		utils.RecordLog(utils.Info, utils.StartTask, "", fmt.Sprintf("TaskId: %s", taskId), 0)
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
	utils.RecordLog(utils.Info, utils.RestartTask, "", "", 0)
	return operations.NewUpdateTaskCreated().WithPayload(&models.OK{
		Message: "ok",
	})
}

func GetTaskCommandOutHandler(params operations.GetTaskCommandOutParams) middleware.Responder {
	out := &models.TaskCommandOut{}
	row := database.DB.QueryRow("select log_id, content from task_command_out where log_id = ?", params.LogID)
	err := row.Scan(&out.LogID, &out.Content)
	if err != nil && err != sql.ErrNoRows {
		return operations.NewGetTaskCommandOutInternalServerError().WithPayload(&models.ErrorModel{
			Code:    http.StatusInternalServerError,
			Message: "Sql Error",
		})
	}
	return operations.NewGetTaskCommandOutOK().WithPayload(out)
}

func UpdateTaskCommandOutHandler(params operations.UpdateTaskCommandOutParams) middleware.Responder {
	sqlStr := "INSERT OR REPLACE INTO task_command_out (log_id,content) VALUES (?,?)"
	database.DB.MustExec(sqlStr, params.TaskCommandOut.LogID, params.TaskCommandOut.Content)
	return operations.NewUpdateTaskCommandOutCreated().WithPayload(&models.OK{
		Message: "ok",
	})
}
