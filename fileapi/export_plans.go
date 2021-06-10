package fileapi

import (
	"bytes"
	"ctgb/database"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	log "github.com/sirupsen/logrus"
)

func WritePlansIntoExcel(plans []*database.PlanModel, tasks []*database.TaskModel) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	// Create a new sheet.
	taskInfo := make(map[int64]string)

	for _, task := range tasks {
		if task.Id.Int64 != 0 {
			if task.Status.String == "completed" {
				taskInfo[task.Id.Int64] = "completed"
			} else {
				taskInfo[task.Id.Int64] = "running"
			}
		}
	}
	collumns := strings.Split("编号，状态，任务类型，PVOB，组件，子目录，仓库地址，迁移方式，目标仓库地址，计划迁移日期，计划切换日期，"+
		"实际迁移日期，实际切换日期，物理子系统，配置库，事业群，项目组，对接人姓名，联系人电话，备注，工程类型，业务用途，影响范围，计划状态", "，")
	az := strings.Split("ABCDEFGHIJKLMNOPQRSTUVWXYZ", "")
	for i, collumn := range collumns {
		f.SetCellValue("Sheet1", az[i]+"1", collumn)
	}
	for i, plan := range plans {
		f.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), plan.ID)
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(i+2), plan.Status)
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(i+2), plan.OriginType)
		f.SetCellValue("Sheet1", "D"+strconv.Itoa(i+2), plan.Pvob)
		f.SetCellValue("Sheet1", "E"+strconv.Itoa(i+2), plan.Component)
		f.SetCellValue("Sheet1", "F"+strconv.Itoa(i+2), plan.Dir)
		f.SetCellValue("Sheet1", "G"+strconv.Itoa(i+2), plan.OriginURL)
		f.SetCellValue("Sheet1", "H"+strconv.Itoa(i+2), plan.TranslateType)
		f.SetCellValue("Sheet1", "I"+strconv.Itoa(i+2), plan.TargetURL)
		f.SetCellValue("Sheet1", "J"+strconv.Itoa(i+2), plan.PlanStartTime)
		f.SetCellValue("Sheet1", "K"+strconv.Itoa(i+2), plan.PlanSwitchTime)
		f.SetCellValue("Sheet1", "L"+strconv.Itoa(i+2), plan.ActualStartTime)
		f.SetCellValue("Sheet1", "M"+strconv.Itoa(i+2), plan.ActualSwitchTime)
		f.SetCellValue("Sheet1", "N"+strconv.Itoa(i+2), plan.Subsystem)
		f.SetCellValue("Sheet1", "O"+strconv.Itoa(i+2), plan.ConfigLib)
		f.SetCellValue("Sheet1", "P"+strconv.Itoa(i+2), plan.Group)
		f.SetCellValue("Sheet1", "Q"+strconv.Itoa(i+2), plan.Team)
		f.SetCellValue("Sheet1", "R"+strconv.Itoa(i+2), plan.Supporter)
		f.SetCellValue("Sheet1", "S"+strconv.Itoa(i+2), plan.SupporterTel)
		f.SetCellValue("Sheet1", "T"+strconv.Itoa(i+2), plan.Tip)
		f.SetCellValue("Sheet1", "U"+strconv.Itoa(i+2), plan.ProjectType)
		f.SetCellValue("Sheet1", "V"+strconv.Itoa(i+2), plan.Purpose)
		f.SetCellValue("Sheet1", "W"+strconv.Itoa(i+2), plan.Effect)
		if plan.TaskID == 0 {
			f.SetCellValue("Sheet1", "X"+strconv.Itoa(i+2), "未迁移")
		} else {
			taskStatus, ok := taskInfo[plan.TaskID]
			if !ok {
				f.SetCellValue("Sheet1", "X"+strconv.Itoa(i+2), "未迁移")
			} else {
				if taskStatus == "completed" {
					f.SetCellValue("Sheet1", "X"+strconv.Itoa(i+2), "已迁移")
				} else {
					f.SetCellValue("Sheet1", "X"+strconv.Itoa(i+2), "迁移中")
				}
			}
		}
	}
	buff, err := f.WriteToBuffer()
	if err != nil {
		log.Error("save excel err:", err)
		return nil, err
	}
	return buff, nil
}

func PlansExportHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now().Format("2006-01-02-15-04-05")
	plans := make([]*database.PlanModel, 0)
	err := database.DB.Select(&plans, "SELECT * FROM plan")
	if err != nil {
		log.Error("export error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	tasks := make([]*database.TaskModel, 0)
	err = database.DB.Select(&tasks, "SELECT * FROM task")
	if err != nil {
		log.Error("export error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	log.Debug("begin plans export at:", now)
	log.Debug("total:", len(plans))

	buff, err := WritePlansIntoExcel(plans, tasks)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	fileData := buff.Bytes()
	if err != nil {
		log.Error("Read File Err:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	} else {
		w.Header().Add("Content-Disposition", "attachment")
		//w.Header().Add("Content-Type", "application/vnd.ms-excel")
		w.Header().Add("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Write(fileData)
	}
}
