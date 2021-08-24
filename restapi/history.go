package restapi

import (
	"ctgb/database"
	"ctgb/models"
	"ctgb/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	"gorm.io/gorm"
)

func GetHistory(params operations.GetCCHistoryParams) middleware.Responder {
	gitName := params.GitName
	id := ""
	if params.ID != nil {
		id = *params.ID
	}
	offset := int(params.Offset)
	limit := int(params.Limit)
	var historyArray []database.History
	var count int64
	if id == "" {
		database.DB.Model(&database.History{}).Where("git_name = ?", gitName).Count(&count)
		database.DB.Where("git_name = ?", gitName).Offset(offset).Limit(limit).Find(&historyArray)
	} else {
		database.DB.Where("git_name = ? AND history_id = ?", gitName, id).Count(&count)
		database.DB.Where("git_name = ? AND history_id = ?", gitName, id).Offset(offset).Limit(limit).Find(&historyArray)
	}
	response := models.CCHistoryInfoModel{
		Count:  count,
		Offset: int64(offset + limit),
		Limit:  int64(limit),
	}
	items := make([]*models.CCHistoryInfoModelItem, 0)
	for _, i := range historyArray {
		item := models.CCHistoryInfoModelItem{
			CreateTime:  i.CreateTime,
			Description: i.Description,
			HistoryType: i.HistoryType,
			ID:          i.HistoryId,
			Name:        i.Name,
			Owner:       i.Owner,
		}
		items = append(items, &item)
	}
	response.InfoItem = items
	return operations.NewGetCCHistoryOK().WithPayload(&response)
}

func CreateHistory(params operations.CreateCCHistoryParams) middleware.Responder {
	gitName := *params.GitName
	action := ""
	if params.Action != nil {
		action = *params.Action
	}
	var tx *gorm.DB
	if action == "delete" {
		tx = database.DB.Where("git_name = ?", gitName).Delete(&database.History{})
	} else {
		id := params.InfoItem.ID
		owner := params.InfoItem.Owner
		historyType := params.InfoItem.HistoryType
		createAt := params.InfoItem.CreateTime
		description := params.InfoItem.Description
		name := params.InfoItem.Name
		tx = database.DB.Create(&database.History{
			GitName:     gitName,
			HistoryId:   id,
			Owner:       owner,
			HistoryType: historyType,
			CreateTime:  createAt,
			Description: description,
			Name:        name,
		})
	}
	if tx.Error == nil {
		return operations.NewCreateCCHistoryCreated()
	} else {
		return operations.NewCreateCCHistoryInternalServerError()
	}
}
