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
	"net/http"
	"strconv"
	"time"
)

func startTask(taskID int64) {
	task := &database.TaskModel{}
	err := database.DB.Get(task, "SELECT cc_password,"+
		" cc_user, component, git_password, git_url, git_user, pvob, include_empty"+
		" FROM task WHERE id = $1", taskID)
	if err != nil {
		fmt.Println(err)
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
		return
	}
	var matchInfo []*models.TaskMatchInfo
	database.DB.Select(&matchInfo, "SELECT git_branch, stream FROM match_info WHERE task_id = $1 ORDER BY id",
		taskID)

	workerTaskModel := struct {
		TaskId       int64
		CcPassword   string
		CcUser       string
		Component    string
		GitPassword  string
		GitURL       string
		GitUser      string
		Pvob         string
		Stream       string
		Branch       string
		IncludeEmpty bool
	}{
		TaskId:       taskID,
		CcPassword:   task.CcPassword,
		CcUser:       task.CcUser,
		Component:    task.Component,
		GitPassword:  task.GitPassword,
		GitURL:       task.GitURL,
		GitUser:      task.GitUser,
		Pvob:         task.Pvob,
		Stream:       matchInfo[0].Stream,
		Branch:       matchInfo[0].GitBranch,
		IncludeEmpty: task.IncludeEmpty,
	}
	workerTaskModelByte, _ := json.Marshal(workerTaskModel)
	req, _ := http.NewRequest(http.MethodPost, "http://"+workerUrl+"/new_task", bytes.NewBuffer(workerTaskModelByte))
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	startTime := time.Now().Format("2006-01-02 15:04:05")
	tx := database.DB.MustBegin()
	tx.MustExec(
		"UPDATE task SET status = 'running', worker_id = $1 WHERE id = $2", worker.Id, taskID,
	)
	tx.MustExec(
		"UPDATE worker SET task_count = task_count + 1 WHERE id = $1", worker.Id,
	)
	tx.MustExec(
		"INSERT INTO task_log (task_id, status, start_time, end_time, duration)"+
			" VALUES($1, 'running', $2, $3, 0)", taskID, startTime, "",
	)
	tx.Commit()
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusCreated {
		fmt.Println(fmt.Errorf("不能发送任务给%d", worker.Id), err)
		return
	}
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
		"git_user, git_password, git_url, status, last_completed_date_time, creator, include_empty)"+
		" VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, '', $10)",
		taskInfo.Pvob, taskInfo.Component, taskInfo.CcUser, taskInfo.CcPassword, taskInfo.GitURL,
		taskInfo.GitUser, taskInfo.GitPassword, taskInfo.GitURL, "init", username, taskInfo.IncludeEmpty)
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
		" cc_user, component, git_password, git_url, git_user, pvob, include_empty"+
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
			" status, include_empty" +
			" FROM task WHERE creator = $1 or 1 = 1 ORDER BY id LIMIT $2 OFFSET $3;"
		queryCount = "SELECT count(id) FROM task;"
	} else {
		query = "SELECT pvob, component, git_url, id, last_completed_date_time," +
			" status, include_empty" +
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
		fmt.Println(err)
		return middleware.Error(404, "没发现任务")
	}
	tx := database.DB.MustBegin()
	if params.TaskLog.Status != "" {
		tx.MustExec("UPDATE task_log SET status = $1, end_time = $2, duration = $3 WHERE log_id = $4",
			taskLogInfo.Status, taskLogInfo.EndTime, taskLogInfo.Duration, taskLog.LogId)
		tx.MustExec("UPDATE task SET status = 'completed', last_completed_date_time = $1 WHERE id = $2",
			taskLogInfo.EndTime, taskId)
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
	tx.Commit()
	return operations.NewUpdateTaskCreated().WithPayload(&models.OK{
		Message: "ok",
	})
}
func RestartTaskHandler(params operations.RestartTaskParams) middleware.Responder {
	//username, verified := utils.Verify(params.Authorization)
	taskId := params.RestartTrigger.ID
	task := &database.TaskModel{}
	database.DB.Get(task, "SELECT status, worker_id FROM task WHERE id = $1", taskId)
	if task.Status == "completed" || task.Status == "init" {
		taskIdInt, _ := strconv.ParseInt(taskId, 10, 64)
		go startTask(taskIdInt)
	}
	return operations.NewUpdateTaskCreated().WithPayload(&models.OK{
		Message: "ok",
	})
}
