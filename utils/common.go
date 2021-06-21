package utils

import (
	"bytes"
	"ctgb/database"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/urlesc"

	"github.com/axgle/mahonia"
)

type CheckTaskInfo struct {
	CCUser     string `json:"cc_user"`
	CCPassword string `json:"cc_password"`
	GitRepoURL string `json:"git_repo_url"`
	ModelType  string `json:"model_type"`
	WorkerURL  string `json:"worker_url,omitempty"`
}

func Iconv(src string, srcCode string, targetCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(targetCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	return string(cdata)
}

func ParseGitURL(user, passwd, gitUrl string) string {
	user = urlesc.QueryEscape(user)
	passwd = urlesc.QueryEscape(passwd)
	if strings.HasPrefix(gitUrl, "http://") {
		gitUrl = strings.Replace(gitUrl, "http://", "", 1)
		gitUrl = "http://" + user + ":" + passwd + "@" + gitUrl
	} else if strings.HasPrefix(gitUrl, "https://") {
		gitUrl = strings.Replace(gitUrl, "https://", "", 1)
		gitUrl = "https://" + user + ":" + passwd + "@" + gitUrl
	}
	return gitUrl
}

func DoCheckInfoReq(taskInfo *CheckTaskInfo) (int, map[string]string) {
	var workerURLs []string
	var errRet = make(map[string]string, 0)
	if taskInfo.WorkerURL == "" {
		database.DB.Select(&workerURLs, "select worker_url from worker where status='running'")
	} else {
		workerURLs = append(workerURLs, taskInfo.WorkerURL)
	}
	for _, workerUrl := range workerURLs {
		taskInfo.WorkerURL = ""
		workerTaskModelByte, _ := json.Marshal(taskInfo)
		req, _ := http.NewRequest(http.MethodPost, "http://"+workerUrl+"/check_info", bytes.NewBuffer(workerTaskModelByte))
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp == nil {
			errRet[workerUrl] = "request fail"
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			body, _ := ioutil.ReadAll(resp.Body)
			errRet[workerUrl] = string(body)
			continue
		}
	}
	if len(errRet) == 0 {
		return http.StatusOK, errRet
	} else {
		return http.StatusInternalServerError, errRet
	}
}
