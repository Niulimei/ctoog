package restapi

import (
	"ctgb/database"
	"ctgb/models"
	"ctgb/restapi/operations"
	"net/http"
	"strconv"
	"time"

	"github.com/go-openapi/runtime/middleware"
	log "github.com/sirupsen/logrus"
)

func CleanDeadWorker() {
	var workers []*database.WorkerModel
	start := time.Now().Add(time.Second * -30).Format("2006-01-02 15:04:05")
	err := database.DB.Select(&workers, "SELECT * FROM worker WHERE register_time < $1 AND status = 'running'", start)
	if err == nil && len(workers) > 0 {
		tx, err := database.DB.Begin()
		if err != nil {
			log.Error("database error:", err)
			return
		}
		for _, worker := range workers {
			tx.Exec("UPDATE worker SET status = 'dead' WHERE id = $1", worker.Id)
		}
		tx.Commit()
	}
}

func init() {
	go func() {
		for {
			CleanDeadWorker()
			time.Sleep(time.Second * 30)
		}
	}()
}

func PingWorkerHandler(params operations.PingWorkerParams) middleware.Responder {
	host := params.WorkerInfo.Host
	port := params.WorkerInfo.Port
	url := host + ":" + strconv.FormatInt(port, 10)
	log.Println("worker url", url)
	worker := &database.WorkerModel{}
	err := database.DB.Get(worker, "SELECT * FROM worker WHERE worker_url = $1", url)
	now := time.Now().Format("2006-01-02 15:04:05")
	if err != nil {
		database.DB.Exec("INSERT INTO worker (worker_url, status, task_count, register_time) "+
			"VALUES ($1, $2, $3, $4)", url, "running", 0, now)
	} else {
		database.DB.Exec("UPDATE worker SET status = 'running', register_time = $1 WHERE id = $2", now, worker.Id)
	}
	return operations.NewPingWorkerCreated().WithPayload(&models.OK{Message: "ok"})
}

func GetWorkerHandler(param operations.GetWorkerParams) middleware.Responder {
	sqlStr := "select * from worker where id=?"
	var workerInfo = &models.WorkerDetail{}
	err := database.DB.Get(workerInfo, sqlStr, param.ID)
	if err != nil {
		return operations.NewGetWorkerInternalServerError().WithPayload(&models.ErrorModel{
			Code:    http.StatusInternalServerError,
			Message: "Sql Error",
		})
	}
	return operations.NewGetWorkerOK().WithPayload(workerInfo)
}

func ListWorkersHandler(param operations.ListWorkersParams) middleware.Responder {
	if !CheckPermission(param.HTTPRequest) {
		return operations.NewListWorkersInternalServerError().WithPayload(&models.ErrorModel{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}
	sqlStr := "select * from worker order by id limit ? offset ?"
	sqlCountStr := "select count(id) from worker"
	var workerInfo []*models.WorkerDetail
	err := database.DB.Select(&workerInfo, sqlStr, param.Limit, param.Offset)
	if err != nil {
		return operations.NewListWorkersInternalServerError().WithPayload(&models.ErrorModel{
			Code:    http.StatusInternalServerError,
			Message: "Sql Error",
		})
	}
	var count int64
	database.DB.Get(&count, sqlCountStr)
	return operations.NewListWorkersOK().WithPayload(&models.WorkerPageInfoModel{
		Count:      count,
		WorkerInfo: workerInfo,
	})
}
