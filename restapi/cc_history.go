package restapi

import (
	"ctgb/models"
	"ctgb/restapi/operations"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"time"
)

func GetCCHistory(params operations.GetCCHistoryParams) middleware.Responder {
	offset := params.Offset
	limit := params.Limit
	_ = params.GitName
	historys := make([]*models.CCHistoryInfoModel, 0)
	var i int64
	for i = 0; i < 10; i++ {
		createTime := time.Now().Add(time.Hour * time.Duration(i) * -1)
		historys = append(historys, &models.CCHistoryInfoModel{
			Key:         i,
			Name:        fmt.Sprintf("测试名称%d", i),
			Owner:       fmt.Sprintf("测试用户%d", i),
			CreateTime:  createTime.Format("2006-01-02 15:04:05"),
			Description: "a commit from " + createTime.Format("2006-01-02 15:04:05"),
			HistoryType: "check in",
			ID:          "random_string",
			Offset:      offset,
			Limit:       limit,
			Count:       10,
		})
	}
	return operations.NewGetCCHistoryOK().WithPayload(historys)
}
