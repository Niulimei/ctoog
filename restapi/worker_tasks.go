package restapi

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/urlesc"
	log "github.com/sirupsen/logrus"
)

func infoServerTaskCompleted(task *Task, server string, cmds []*exec.Cmd, tmpCmdOutFile string, endSignal chan struct{}) {
	data := payload{
		Logid:  strconv.FormatInt(task.TaskLogId, 10),
		Status: "completed",
	}
	start := time.Now()
	for _, cmd := range cmds {
		data.Starttime = start.Format("2006-01-02 15:04:05")
		log.Debug("start cmd:", cmd.String())
		err := sendCommandOut(server, cmd, task, tmpCmdOutFile)
		if err != nil {
			log.Error("cmd err:", err)
			data.Status = "failed"
			break
		}
	}
	endSignal <- struct{}{}
	end := time.Now()
	data.Endtime = end.Format("2006-01-02 15:04:05")
	duration := end.Sub(start).Seconds()
	d := strconv.FormatInt(int64(duration), 10)
	data.Duration = d
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		log.Error(err)
		return
	}
	log.Debugf("playlod: %+v\n", data)
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("PUT",
		fmt.Sprintf("http://%s/api/tasks/%d?start=false", server, task.TaskId), body)
	if err != nil {
		log.Error("create request error:", err)
		return
		// handle err
	}
	doSend(req)
}

func sendCommandOut(server string, cmd *exec.Cmd, task *Task, tmpCmdOutFile string) error {
	var stop = make(chan struct{})
	data := &commandOut{
		Logid: task.TaskLogId,
	}
	//stdout, err := cmd.StdoutPipe()
	//if err != nil {
	//	log.Println(err)
	//	return err
	//}
	//stderr, err := cmd.StderrPipe()
	//if err != nil {
	//	log.Println(err)
	//	return err
	//}

	if err := cmd.Start(); err != nil {
		log.Errorln(err)
		return err
	}

	//s := bufio.NewScanner(io.MultiReader(stdout, stderr))
	tk := time.NewTicker(time.Second * 1)
	var logContent []string
	go func(in *[]string) {
		for {
			select {
			case <-tk.C:
				data.Content = strings.Join(*in, "\n")
				sender(server, data)
			case <-stop:
				return
			}
		}
	}(&logContent)
	go readCommandOut(tmpCmdOutFile, &logContent)
	//for s.Scan() {
	//	tmp = append(tmp, s.Text())
	//	//tmp = append(tmp, utils.Iconv(s.Text(), "gbk", "utf8"))
	//}

	err := cmd.Wait()
	log.Errorln(err)
	go func() {
		var count, lastTimeLen int
		lastTimeLen = len(logContent)
		for {
			if len(logContent) == lastTimeLen {
				count++
				if count >= 3 {
					break
				}
			}
			lastTimeLen = len(logContent)
			time.Sleep(time.Second)
		}
		stop <- struct{}{}
		os.RemoveAll(tmpCmdOutFile)
	}()
	return err
}

func sender(server string, data *commandOut) {
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		log.Error(err)
	}
	log.Debugf("playlod: %+v\n", data)
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST",
		fmt.Sprintf("http://%s/api/tasks/cmdout/%d", server, data.Logid), body)
	if err != nil {
		return
		// handle err
	}
	doSend(req)
}

func doSend(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", adminJwtToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.Body == nil {
		// handle err
		time.Sleep(time.Second * 3)
		resp, err = http.DefaultClient.Do(req)
		if err != nil || resp.Body == nil {
			log.Error("send to server err:", err)
			return
		}
	}
	log.Debug("info server success")
	resp.Body.Close()
}

func PingServer(host string, port int) {
	defer func() {
		if ret := recover(); ret != nil {
			fmt.Printf("Recover From Panic. %v\n", ret)
		}
	}()

	data := Payload{
		Host: host,
		Port: port,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return
		// handle err
	}
	for {
		body := bytes.NewReader(payloadBytes)
		req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/api/workers", ServerFlag), body)
		if err != nil {
			return
			// handle err
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", adminJwtToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Error(err)
		} else {
			if resp.Body != nil {
				resp.Body.Close()
			}
		}
		time.Sleep(time.Second * 10)
	}
}

func DeleteWorkerTaskCacheHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("read task error:", err)
		return
	}
	if r.Body != nil {
		defer r.Body.Close()
	}
	workerTaskModel := WorkerTaskDelInfo{}
	err = json.Unmarshal(body, &workerTaskModel)
	if err != nil {
		log.Error(err)
		w.WriteHeader(500)
		w.Write([]byte("Json marshal fail"))
		return
	}
	delCache(w, workerTaskModel)
}

func delCache(w http.ResponseWriter, workerTaskModel WorkerTaskDelInfo) {
	cwd, _ := os.Getwd()
	var checkCacheCmdStr, cleanCacheCmdStr string
	switch workerTaskModel.ModelType {
	case "clearcase":
		checkCacheCmdStr = fmt.Sprintf(`/usr/bin/bash %s/script/cc2git/checkCache.sh %d`, cwd, workerTaskModel.TaskId)
		cleanCacheCmdStr = fmt.Sprintf(`echo %s | su - %s -c "/usr/bin/bash %s/script/cc2git/cleanCache.sh %d %s"`,
			workerTaskModel.CcPassword, workerTaskModel.CcUser, cwd, workerTaskModel.TaskId, workerTaskModel.Exception)
	case "svn":
		checkCacheCmdStr = fmt.Sprintf(`/usr/bin/bash %s/script/svn2git/checkCache.sh %d`, cwd, workerTaskModel.TaskId)
		cleanCacheCmdStr = fmt.Sprintf(`/usr/bin/bash %s/script/svn2git/cleanCache.sh %d`, cwd, workerTaskModel.TaskId)
	default:
		w.WriteHeader(500)
		w.Write([]byte("Not Support"))
		return
	}

	//目录不存在直接返回成功
	log.Infoln(checkCacheCmdStr)
	cmd := exec.Command("/bin/bash", "-c", checkCacheCmdStr)
	_, err := cmd.Output()
	if err != nil {
		log.Errorln(err.Error())
		w.WriteHeader(200)
		w.Write([]byte("success"))
		return
	}

	log.Infoln(cleanCacheCmdStr)
	cmd = exec.Command("/bin/bash", "-c", cleanCacheCmdStr)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	cmd.Start()
	s := bufio.NewScanner(io.MultiReader(stdout, stderr))
	var tmp []string
	for s.Scan() {
		tmp = append(tmp, s.Text())
	}
	err = cmd.Wait()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(strings.Join(tmp, "\n")))
		return
	}
	w.WriteHeader(200)
	w.Write([]byte("success"))
}

func WorkerTaskHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("read task error:", err)
		return
	}
	if r.Body != nil {
		defer r.Body.Close()
	}
	workerTaskModel := Task{}
	err = json.Unmarshal(body, &workerTaskModel)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Json marshal fail"))
		return
	}
	gitUrl := parseGitURL(workerTaskModel.GitUser, workerTaskModel.GitPassword, workerTaskModel.GitURL)
	switch workerTaskModel.ModelType {
	case "clearcase":
		cc2Git(workerTaskModel, gitUrl)
	case "svn":
		svn2Git(workerTaskModel, gitUrl)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Not Support"))
	}

	w.WriteHeader(201)
	w.Write([]byte("bye"))
}

func readCommandOut(fileName string, container *[]string) {
	var stopRead = make(chan struct{})
	file, err := os.Open(fileName)
	if err != nil {
		log.Errorf("Open file fail:%v", err)
		return
	}
	defer file.Close()
	go func() {
		for {
			_, err := os.Lstat(fileName)
			if os.IsNotExist(err) {
				stopRead <- struct{}{}
				return
			}
			time.Sleep(time.Second)
		}
	}()
	reader := bufio.NewReader(file)
	var tick = time.NewTicker(100 * time.Millisecond)
	func() {
		for {
			select {
			case <-tick.C:
				line, err := reader.ReadString('\n')
				if err != nil {
					if err != io.EOF {
						return
					}
				}
				if line != "" {
					*container = append(*container, line)
				}
			case <-stopRead:
				return
			}
		}
	}()
}

func cc2Git(workerTaskModel Task, gitUrl string) {
	var endSignal = make(chan struct{})
	cwd, _ := os.Getwd()
	var cmds []*exec.Cmd
	tmpCmdOutFile := fmt.Sprintf("%s/tmpCmdOut/%d_%d.log", cwd, workerTaskModel.TaskId, workerTaskModel.TaskLogId)
	exec.Command("/bin/bash", "-c", fmt.Sprintf("mkdir -p %s/tmpCmdOut;touch %s/tmpCmdOut/%d_%d.log", cwd, cwd, workerTaskModel.TaskId, workerTaskModel.TaskLogId)).Output()
	for _, match := range workerTaskModel.Matches {
		cmdStr := fmt.Sprintf(`echo %s | su - %s -c '/usr/bin/bash %s/script/cc2git/cc2git.sh "%s" "%s" "%s" "%s" "%s" "%d" "%t" "%s" "%s" "%s" "%s"' &>>%s`,
			workerTaskModel.CcPassword, workerTaskModel.CcUser, cwd, workerTaskModel.Pvob, workerTaskModel.Component,
			match.Stream, gitUrl, match.Branch, workerTaskModel.TaskId,
			workerTaskModel.IncludeEmpty, workerTaskModel.GitUser, workerTaskModel.GitEmail, workerTaskModel.Keep,
			strings.ReplaceAll(workerTaskModel.Gitignore, " ", ""), tmpCmdOutFile)
		log.Infoln(cmdStr)
		cmd := exec.Command("/bin/bash", "-c", cmdStr)
		cmds = append(cmds, cmd)
	}
	go infoServerTaskCompleted(&workerTaskModel, ServerFlag, cmds, tmpCmdOutFile, endSignal)
}

func geneUsersFile(workerTaskModel Task) string {
	cwd, _ := os.Getwd()
	var buffer bytes.Buffer
	for _, pi := range workerTaskModel.NamePair {
		buffer.WriteString(fmt.Sprintf("%s = %s <%s>\n", pi.SnvUserName, pi.GitUserName, pi.GitEmail))
	}
	fp := filepath.Join(cwd, "tmpCmdOut", filepath.Base(workerTaskModel.SvnURL)+"_"+strconv.Itoa(int(workerTaskModel.TaskId))+".txt")
	ioutil.WriteFile(fp, []byte(buffer.String()), 0644)
	return fp
}

func svn2Git(workerTaskModel Task, gitUrl string) int {
	var endSignal = make(chan struct{})
	cwd, _ := os.Getwd()
	var cmds []*exec.Cmd
	tmpCmdOutFile := fmt.Sprintf("%s/tmpCmdOut/%d_%d.log", cwd, workerTaskModel.TaskId, workerTaskModel.TaskLogId)
	exec.Command("/bin/bash", "-c", fmt.Sprintf("mkdir -p %s/tmpCmdOut;touch %s/tmpCmdOut/%d_%d.log", cwd, cwd, workerTaskModel.TaskId, workerTaskModel.TaskLogId)).Output()
	userFile := geneUsersFile(workerTaskModel)
	cmdStr := fmt.Sprintf(`/usr/bin/bash %s/script/svn2git/svn2git.sh "%s" "%s" "%d" "%t" "%s" "%s" "%s" "%s" "%s" &> %s`,
		cwd, workerTaskModel.SvnURL, gitUrl, workerTaskModel.TaskId,
		workerTaskModel.IncludeEmpty, workerTaskModel.GitUser, workerTaskModel.GitEmail,
		workerTaskModel.Keep, userFile, strings.ReplaceAll(workerTaskModel.Gitignore, " ", ""),
		tmpCmdOutFile)
	log.Infoln(cmdStr)
	cmd := exec.Command("/bin/bash", "-c", cmdStr)
	cmds = append(cmds, cmd)
	go infoServerTaskCompleted(&workerTaskModel, ServerFlag, cmds, tmpCmdOutFile, endSignal)
	go func() {
		<-endSignal
		os.RemoveAll(userFile)
	}()
	return http.StatusOK
}

func parseGitURL(user, passwd, gitUrl string) string {
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
