package restapi

import (
	"ctgb/database"
	"ctgb/models"
	"ctgb/restapi/operations"
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
