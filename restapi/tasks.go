package restapi

import (
	"ctgb/database"
	"ctgb/models"
	"ctgb/restapi/operations"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"strings"
)

func CreateTaskHandler(params operations.CreateTaskParams) middleware.Responder {
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
		"git_user, git_password, git_repo, status, last_completed_date_time) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, '')",
		taskInfo.Pvob, taskInfo.Component, taskInfo.CcUser, taskInfo.CcPassword, taskInfo.GitURL,
		taskInfo.GitUser, taskInfo.GitPassword, repo, "init")
	taskId, err := r.LastInsertId()
	if err != nil {
		return operations.NewCreateTaskInternalServerError().WithPayload(
			&models.ErrorModel{Message: fmt.Sprintf("Insert into db error: %+v", err), Code: 500})
	}
	tx, _ := database.DB.Begin()
	for _, match := range taskInfo.MatchInfo {
		tx.Exec("INSERT INTO " +
			"match_info (task_id, stream, git_branch) " +
			"VALUES($1, $2, $3)",
			taskId, match.Stream, match.GitBranch)
	}
	tx.Commit()
	return operations.NewCreateTaskCreated().WithPayload(&models.OK{Message: "ok"})
}

func GetTaskHandler(params operations.GetTaskParams) middleware.Responder {
	taskID := params.ID
	task := &models.TaskModel{}
	database.DB.Get(task,"SELECT cc_password," +
		" cc_user, component, git_password, git_url, git_user, pvob" +
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
	var tasks []*models.TaskInfoModel
	database.DB.Select(&tasks, "SELECT pvob, component, git_repo, id, last_completed_date_time," +
		" status" +
		" FROM task;")
	return operations.NewListTaskOK().WithPayload(tasks)
}
