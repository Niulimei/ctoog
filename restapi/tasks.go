package restapi

import (
	"ctgb/database"
	"ctgb/models"
	"ctgb/restapi/operations"
	"ctgb/utils"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"strings"
)

func startTask(taskID string) error {
	task := &database.TaskModel{}
	err := database.DB.Get(task, "SELECT cc_password,"+
		" cc_user, component, git_password, git_url, git_user, pvob"+
		" FROM task WHERE id = $1", taskID)
	if err != nil {
		return nil
	}
	worker := &database.WorkerModel{}
	if task.WorkerId != 0 {
		database.DB.Get(worker, "SELECT * FROM worker WHERE id = $1", task.WorkerId)
	} else {
		database.DB.Select(worker, "SELECT * FROM worker ORDER BY task_count DESC limit 1")
	}
	_ = worker.WorkerUrl
	return nil
}

func CreateTaskHandler(params operations.CreateTaskParams) middleware.Responder {
	userToken := params.Authorization
	username, verified := utils.Verify(userToken)
	if !verified {
		return middleware.Error(401, "鉴权失败")
	}
	fmt.Println(username)
	taskInfo := params.TaskInfo
	urls := strings.Split(taskInfo.GitURL, "/")
	var repo string
	if len(urls) > 1 {
		repo = urls[len(urls)-2]
		repo = strings.Replace(repo, ".git", "", 1)
	} else {
		repo = ""
	}
	r := database.DB.MustExec("INSERT INTO task (pvob, component, cc_user, cc_password, git_url,"+
		"git_user, git_password, git_repo, status, last_completed_date_time, creator)"+
		" VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, '', $10)",
		taskInfo.Pvob, taskInfo.Component, taskInfo.CcUser, taskInfo.CcPassword, taskInfo.GitURL,
		taskInfo.GitUser, taskInfo.GitPassword, repo, "init", username)
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
	return operations.NewCreateTaskCreated().WithPayload(&models.OK{Message: "ok"})
}

func GetTaskHandler(params operations.GetTaskParams) middleware.Responder {
	taskID := params.ID
	task := &models.TaskModel{}
	database.DB.Get(task, "SELECT cc_password,"+
		" cc_user, component, git_password, git_url, git_user, pvob"+
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
		return middleware.Error(401, "鉴权失败")
	}
	var query string
	if username == "admin" {
		query = "SELECT pvob, component, git_repo, id, last_completed_date_time," +
			" status" +
			" FROM task WHERE creator = $1 or 1 = 1 ORDER BY id OFFSET $2 LIMIT 10;"
	} else {
		query = "SELECT pvob, component, git_repo, id, last_completed_date_time," +
			" status" +
			" FROM task WHERE creator = $1 ORDER BY id OFFSET $2 LIMIT 10;"
	}
	offset := params.Offset
	var tasks []*models.TaskInfoModel
	var count int64
	database.DB.Get(count, query, username, offset)
	tasksPage := &models.TaskPageInfoModel{}
	tasksPage.TaskInfo = tasks
	return operations.NewListTaskOK().WithPayload(tasksPage)
}

func UpdateTaskHandler(params operations.UpdateTaskParams) middleware.Responder {
	//username, verified := utils.Verify(params.Authorization)
	taskId := params.ID
	taskLogInfo := params.TaskLog
	taskLog := &database.TaskLog{}
	err := database.DB.Get(taskLog, "SELECT * FROM task_log WHERE task_id = $1 AND status = 'running'", taskId)
	if err != nil {
		return middleware.Error(404, "没发现任务")
	}
	database.DB.Exec("UPDATE task SET status = $1, end_time = $2, duration = $3 WHERE id = $4",
		taskLogInfo.Status, taskLogInfo.EndTime, taskLogInfo.Duration, taskLog.LogId)
	return operations.NewUpdateTaskCreated().WithPayload(&models.OK{
		Message: "ok",
	})
}
