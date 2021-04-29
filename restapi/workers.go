package restapi

import (
	"ctgb/database"
	"ctgb/models"
	"ctgb/restapi/operations"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"strconv"
)

func PingWorkerHandler(params operations.PingWorkerParams) middleware.Responder {
	host := params.WorkerInfo.Host
	port := params.WorkerInfo.Port
	url := host + ":" + strconv.FormatInt(port, 10)
	fmt.Println("worker url", url)
	worker := &database.WorkerModel{}
	err := database.DB.Get(worker, "SELECT * FROM worker WHERE worker_url = $1", url)
	fmt.Println(err)
	if err != nil {
		database.DB.Exec("INSERT INTO worker (worker_url, status, task_count) "+
			"VALUES ($1, $2, $3)", url, "running", 0)
	} else {
		database.DB.Exec("UPDATE worker SET status = 'running' WHERE id = $1", worker.Id)
	}
	return operations.NewPingWorkerCreated().WithPayload(&models.OK{Message: "ok"})
}
