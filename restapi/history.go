package restapi

import (
	"ctgb/database"
	"ctgb/models"
	"ctgb/restapi/operations"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"gorm.io/gorm"
	"strings"
	"time"
)

func GetHistory(params operations.GetCCHistoryParams) middleware.Responder {
	gitName := params.GitName
	id := ""
	if params.ID != nil {
		id = *params.ID
		id = strings.Trim(id, " ")
	}
	id = strings.Trim(id, " ")
	offset := int(params.Offset)
	limit := int(params.Limit)
	var historyArray []database.History
	var count int64
	if id == "" {
		database.DB.Model(&database.History{}).Where("git_name = ?", gitName).Count(&count)
		database.DB.Where("git_name = ?", gitName).Order("create_time desc").Offset(offset).Limit(limit).Find(&historyArray)
	} else {
		database.DB.Model(&database.History{}).Where("git_name = ? AND history_id like ?", gitName, "%"+id+"%").Count(&count)
		database.DB.Where("git_name = ? AND history_id like ?", gitName, "%"+id+"%").Offset(offset).Limit(limit).Find(&historyArray)
	}
	response := models.CCHistoryInfoModel{
		Count:  count,
		Offset: int64(offset + limit),
		Limit:  int64(limit),
	}
	items := make([]*models.CCHistoryInfoModelItem, 0)
	for _, i := range historyArray {
		//20140808T10:25:39+08:00
		createTime, err := time.Parse("20060102T15:04:05-07:00", i.CreateTime)
		if err == nil {
			i.CreateTime = createTime.Format("2006/1/2 15:04:05")
		}
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
	token := params.AuthToken
	if token != "052353be3300b083408c6bb6deb2ab67" {
		return middleware.Error(403, "token invalid")
	}
	gitName := *params.GitName
	action := ""
	if params.Action != nil {
		action = *params.Action
	}
	var tx *gorm.DB
	if action == "delete" {
		tx = database.DB.Where("git_name = ?", gitName).Delete(&database.History{})
	} else if action == "update" {
		oldGitName := ""
		if params.OldGitName != nil {
			oldGitName = *params.OldGitName
		}
		tx = database.DB.Model(&database.History{}).Where("git_name = ?", oldGitName).Update("git_name", gitName)
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
		return operations.NewCreateCCHistoryCreated().WithPayload(&models.OK{
			Message: "ok",
		})
	} else {
		return operations.NewCreateCCHistoryInternalServerError().WithPayload(&models.ErrorModel{Message: tx.Error.Error()})
	}
}

func GetHistoryId(params operations.SearchCCHistoryParams) middleware.Responder {
	id := params.ID
	id = strings.Trim(id, " ")
	limit := int(params.Limit)
	gitName := params.GitName
	var ids []string
	var count int64
	database.DB.Model(&database.History{}).Where("git_name = ? AND history_id like ?", gitName, "%"+id+"%").Count(&count)
	tx := database.DB.Table("histories").Select("history_id").Where("git_name = ? AND history_id like ?", gitName, "%"+id+"%").Limit(limit).Find(&ids)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return operations.NewSearchCCHistoryInternalServerError()
	} else {
		return operations.NewSearchCCHistoryOK().WithPayload(ids)
	}
}
