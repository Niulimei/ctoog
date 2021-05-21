package fileapi

import (
	"bytes"
	"ctgb/database"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func WritePlansIntoExcel(plans []*database.PlanModel, dest string) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	// Create a new sheet.
	collumns := strings.Split("编号，状态，任务类型，PVOB，组件，子目录，仓库地址，迁移方式，目标仓库地址，计划迁移日期，计划切换日期，"+
		"实际迁移日期，实际切换日期，物理子系统，配置库，事业群，项目组，对接人姓名，联系人电话，备注，工程类型，业务用途，影响范围", "，")
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
	}
	buff, err := f.WriteToBuffer()
	if err != nil {
		log.Error("save excel err:", err)
		return nil, err
	}
	return buff, nil
}

func PlansExportHandler(w http.ResponseWriter, r *http.Request) {
	pwd, _ := os.Getwd()
	now := time.Now().Format("2006-01-02-15-04-05")
	des := pwd + string(os.PathSeparator) + "planDataExport-" + now + ".xlsx"
	plans := make([]*database.PlanModel, 0)
	err := database.DB.Select(&plans, "SELECT * FROM plan")
	if err != nil {
		log.Error("export error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	log.Debug("begin plans export at:", now)
	log.Debug("total:", len(plans))

	buff, err := WritePlansIntoExcel(plans, des)
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
		log.Println("Send File:", des)
		w.Header().Add("Content-Disposition", "attachment")
		//w.Header().Add("Content-Type", "application/vnd.ms-excel")
		w.Header().Add("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Write(fileData)
	}
}
