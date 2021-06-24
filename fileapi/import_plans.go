package fileapi

import (
	"bytes"
	"ctgb/database"
	"ctgb/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	log "github.com/sirupsen/logrus"
)

func checkData(rows [][]string) (bool, string) {
	var checkColumn = true
	for index, row := range rows {
		if len(row) < 34 {
			break
		}
		taskType := row[2]
		if taskType != "ClearCase" && taskType != "私服" && taskType != "ICDP(Gerrit)" {
			log.Error("taskType not support: ", taskType)
			return false, fmt.Sprintf("第%d行 任务类型错误", index+2)
		}
		if taskType != "ClearCase" {
			continue
		}
		var needColumn = []int{2, 3, 4, 6, 7, 10, 11, 12, 15, 16, 19}
		for _, i := range needColumn {
			if row[i] == "" {
				log.Info("row index is blank:", index)
				checkColumn = false
				return false, fmt.Sprintf("第%d行第%d列数据不全", index+2, i)
			}
		}
		var includeEmpty bool
		if row[8] == "是" || row[8] == "Yes" || row[8] == "yes" || row[8] == "Y" {
			includeEmpty = true
		} else {
			includeEmpty = false
		}
		keep := row[9]
		if includeEmpty && keep == "" {
			return false, fmt.Sprintf("第%d行第%d列数据不全", index+2, 10)
		}
		ccUser := row[10]
		ccPassword := row[11]
		gitUser := row[15]
		gitPassword := row[16]
		gitUrl := row[19]
		checkInfos := &utils.CheckTaskInfo{
			CCUser:     ccUser,
			CCPassword: ccPassword,
			GitRepoURL: utils.ParseGitURL(gitUser, gitPassword, gitUrl),
			ModelType:  strings.ToLower(taskType),
			SvnURL:     "",
		}
		code, errRet := utils.DoCheckInfoReq(checkInfos)
		if code != http.StatusOK {
			ret, _ := json.Marshal(errRet)
			return false, fmt.Sprintf("第%d行 用户或仓库信息错误: %v", index+2, string(ret))
		}
	}
	return checkColumn, ""
}

func PlansImportHandler(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)

	username, _ := utils.Verify(r.Header.Get("Authorization"))

	// Get handler for filename, size and headers
	file, _, err := r.FormFile("uploadFile")
	if err != nil {
		log.Error("Error Retrieving the File")
		log.Error(err)
		return
	}

	defer file.Close()
	buf := bytes.NewBuffer(nil)

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(buf, file); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	excel, err := excelize.OpenReader(buf)
	if err != nil {
		log.Error("open upload file err:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rows, err := excel.GetRows("Sheet1")
	if err != nil {
		log.Error("open upload file err:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rows = rows[1:]
	check, message := checkData(rows)
	if !check {
		log.Error(message)
		http.Error(w, message, http.StatusBadRequest)
		return
	}
	for _, row := range rows {
		if len(row) < 34 {
			log.Info("不足34，跳过:", row)
			continue
		}
		taskType := row[2]
		if taskType != "ClearCase" && taskType != "ICDP(Gerrit)" && taskType != "私服" {
			continue
		}
		log.Info("taskType", taskType)
		pvob := row[3]
		component := row[4]
		dir := row[5]
		stream := row[6]
		branch := row[7]
		var includeEmpty bool
		if row[8] == "是" || row[8] == "Yes" || row[8] == "yes" || row[8] == "Y" {
			includeEmpty = true
		} else {
			includeEmpty = false
		}
		keep := row[9]
		ccUser := row[10]
		ccPassword := row[11]
		workType := row[12]
		gitUser := row[15]
		gitPassword := row[16]
		gitignore := row[18]
		gitEmail := gitUser + "@ccbft.com"
		gitUrl := row[19]
		planStartTime := row[20]
		planSwitchTime := row[21]
		actualStartTime := row[22]
		actualSwitchTime := row[23]
		subsystem := row[24]
		configLib := row[25]
		group := row[26]
		team := row[27]
		supporterName := row[28]
		supporterTel := row[29]
		tip := row[30]
		projectType := row[31]
		purpose := row[32]
		effect := row[33]
		var taskId int64
		var status = "未迁移"

		if taskType == "ClearCase" {
			status = "迁移中"
			r := database.DB.MustExec("INSERT INTO task (pvob, component, cc_user, cc_password, git_url,"+
				"git_user, git_password, status, last_completed_date_time, creator, include_empty, git_email, dir, keep, worker_id, model_type, gitignore)"+
				" VALUES ($1, $2, $3, $4, $5, $6, $7, $8, '', $9, $10, $11, $12, $13, 0, 'clearcase', $14)",
				pvob, component, ccUser, ccPassword, gitUrl, gitUser, gitPassword, "init", username,
				includeEmpty, gitEmail, dir, keep, gitignore)
			taskId, err = r.LastInsertId()
			log.Info("Insert task:", taskId, " ", err)
			database.DB.Exec("INSERT INTO "+
				"match_info (task_id, stream, git_branch) "+
				"VALUES($1, $2, $3)",
				taskId, stream, branch)
		}
		var planColumns = []string{"id", "status", "origin_type", "pvob", "component", "dir", "origin_url", "translate_type", "target_url", "subsystem", "config_lib", "business_group", "team", "supporter", "supporter_tel", "creator", "tip", "project_type", "purpose", "plan_start_time", "plan_switch_time", "actual_start_time", "actual_switch_time", "effect", "task_id", "extra1", "extra2", "extra3"}
		var ph []string
		for i := 1; i <= len(planColumns[1:]); i++ {
			ph = append(ph, "?")
		}
		sqlStr := fmt.Sprintf("INSERT INTO plan (%s) VALUES (%s)",
			strings.Join(planColumns[1:], ","), strings.Join(ph, ","))
		_, err = database.DB.Exec(sqlStr,
			status,
			taskType,
			pvob,
			component,
			dir,
			"",
			workType,
			gitUrl,
			subsystem,
			configLib,
			group,
			team,
			supporterName,
			supporterTel,
			username,
			tip,
			projectType,
			purpose,
			planStartTime,
			planSwitchTime,
			actualStartTime,
			actualSwitchTime,
			effect,
			taskId,
			"",
			"",
			"",
		)
		if err != nil {
			log.Error("insert plan err:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	fmt.Fprintf(w, "Successfully Uploaded File\n")
}
