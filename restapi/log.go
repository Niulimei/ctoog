package restapi

import (
	"ctgb/database"
	"ctgb/models"
	"ctgb/restapi/operations"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-openapi/runtime/middleware"
)

func ListLogsHandler(param operations.ListLogsParams) middleware.Responder {
	var rows *sql.Rows
	var err error
	whereStr := buildWhereSqlStr(param)
	if whereStr == "" {
		rows, err = database.DB.Query("SELECT count(1) over() AS total_rows,"+
			"time,level,user,action,position,message,errcode FROM log ORDER BY time desc LIMIT ? OFFSET ?",
			param.Params.Limit, param.Params.Offset)
	} else {
		rows, err = database.DB.Query(fmt.Sprintf("SELECT count(1) over() AS total_rows,"+
			"time,level,user,action,position,message,errcode FROM log %s ORDER BY time desc LIMIT %d OFFSET %d",
			whereStr, param.Params.Limit, param.Params.Offset))
	}
	if err != nil && err != sql.ErrNoRows {
		return operations.NewListLogsInternalServerError().WithPayload(&models.ErrorModel{
			Code:    http.StatusInternalServerError,
			Message: "Sql Error",
		})
	}
	defer rows.Close()
	var logs []*models.LogInfoModel
	var count int64
	for rows.Next() {
		tmp := &models.LogInfoModel{}
		if err := rows.Scan(&count, &tmp.Time, &tmp.Level,
			&tmp.User, &tmp.Action, &tmp.Position, &tmp.Message, &tmp.Errcode); err != nil {
			return operations.NewListLogsInternalServerError().WithPayload(&models.ErrorModel{
				Code:    http.StatusInternalServerError,
				Message: "Sql Error",
			})
		}
		logs = append(logs, tmp)
	}
	return operations.NewListLogsCreated().WithPayload(&models.LogPageInfoModel{
		Count:   count,
		Limit:   param.Params.Limit,
		Offset:  param.Params.Offset,
		LogInfo: logs,
	})
}

func buildWhereSqlStr(param operations.ListLogsParams) (ret string) {
	var tmp []string
	if param.Params.Action != "" {
		tmp = append(tmp, " action='"+param.Params.Action+"' ")
	}
	if param.Params.Level != "" {
		tmp = append(tmp, " level='"+param.Params.Level+"' ")
	}
	if param.Params.User != "" {
		tmp = append(tmp, " user='"+param.Params.User+"' ")
	}
	if len(tmp) != 0 {
		ret = "where " + strings.Join(tmp, "and")
	}
	return
}
