package main

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
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/urlesc"
	log "github.com/sirupsen/logrus"
)

func infoServerTaskCompleted(task *Task, server string, cmds []*exec.Cmd) {
	data := payload{
		Logid:  strconv.FormatInt(task.TaskLogId, 10),
		Status: "completed",
	}
	start := time.Now()
	for _, cmd := range cmds {
		data.Starttime = start.Format("2006-01-02 15:04:05")
		log.Debug("start cmd:", cmd.String())
		err := sendCommandOut(server, cmd, task)
		if err != nil {
			log.Error("cmd err:", err)
			data.Status = "failed"
			break
		}
	}
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
	log.Printf("playlod: %+v\n", data)
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

func sendCommandOut(server string, cmd *exec.Cmd, task *Task) error {
	data := &commandOut{
		Logid: task.TaskLogId,
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Println(err)
		return err
	}

	if err := cmd.Start(); err != nil {
		log.Println(err)
		return err
	}

	s := bufio.NewScanner(io.MultiReader(stdout, stderr))
	tk := time.NewTicker(time.Second * 1)
	var tmp []string
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
	}(&tmp)
	for s.Scan() {
		tmp = append(tmp, s.Text())
		//tmp = append(tmp, utils.Iconv(s.Text(), "gbk", "utf8"))
	}
	time.Sleep(time.Second * 2)
	stop <- struct{}{}

	if err := cmd.Wait(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func sender(server string, data *commandOut) {
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		log.Error(err)
	}
	log.Printf("playlod: %+v\n", data)
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
	log.Info("info server success")
	resp.Body.Close()
}

func pingServer(host string, port int) {
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
		req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/api/workers", serverFlag), body)
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

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("read task error:", err)
		return
	}
	if r.Body != nil {
		defer r.Body.Close()
	}
	workerTaskModel := TaskDelInfo{}
	if err := json.Unmarshal(body, &workerTaskModel); err == nil {
		cwd, _ := os.Getwd()

		//目录不存在直接返回成功
		cmdStr := fmt.Sprintf(`/usr/bin/bash %s/checkCache.sh %d`, cwd, workerTaskModel.TaskId)
		log.Infoln(cmdStr)
		cmd := exec.Command("/bin/bash", "-c", cmdStr)
		_, err := cmd.Output()
		if err != nil {
			log.Errorln(err.Error())
			w.WriteHeader(200)
			w.Write([]byte("success"))
			return
		}

		cmdStr = fmt.Sprintf(`echo %s | su - %s -c "/usr/bin/bash %s/cleanCache.sh %d %s"`,
			workerTaskModel.CcPassword, workerTaskModel.CcUser, cwd, workerTaskModel.TaskId, workerTaskModel.Exception)
		log.Infoln(cmdStr)
		cmd = exec.Command("/bin/bash", "-c", cmdStr)
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
	} else {
		log.Error(err)
		w.WriteHeader(500)
		w.Write([]byte("fail"))
		return
	}
	w.WriteHeader(200)
	w.Write([]byte("success"))
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("read task error:", err)
		return
	}
	if r.Body != nil {
		defer r.Body.Close()
	}
	workerTaskModel := Task{}
	code := cc2Git(body, workerTaskModel)
	if code != http.StatusOK {
		w.WriteHeader(code)
		w.Write([]byte("fail"))
	}
	w.WriteHeader(201)
	w.Write([]byte("bye"))
}

func cc2Git(body []byte, workerTaskModel Task) int {
	if err := json.Unmarshal(body, &workerTaskModel); err == nil {
		gitUrl := parseGitURL(workerTaskModel.GitUser, workerTaskModel.GitPassword, workerTaskModel.GitURL)
		cwd, _ := os.Getwd()
		var cmds []*exec.Cmd
		for _, match := range workerTaskModel.Matches {
			cmd := exec.Command("/bin/bash", "-c",
				fmt.Sprintf(`echo %s | su - %s -c "/usr/bin/bash %s/script/cc2git/cc2git.sh %s %s %s %s %s %d %t %s %s %s"`,
					workerTaskModel.CcPassword, workerTaskModel.CcUser, cwd, workerTaskModel.Pvob, workerTaskModel.Component,
					match.Stream, gitUrl, match.Branch, workerTaskModel.TaskId,
					workerTaskModel.IncludeEmpty, workerTaskModel.GitUser, workerTaskModel.GitEmail, workerTaskModel.Keep))
			cmds = append(cmds, cmd)
		}
		go infoServerTaskCompleted(&workerTaskModel, serverFlag, cmds)
		return http.StatusOK
	} else {
		log.Error(err)
		return http.StatusInternalServerError
	}
}

func svn2Git(body []byte, workerTaskModel Task) int {
	if err := json.Unmarshal(body, &workerTaskModel); err == nil {
		gitUrl := parseGitURL(workerTaskModel.GitUser, workerTaskModel.GitPassword, workerTaskModel.GitURL)
		cwd, _ := os.Getwd()
		var cmds []*exec.Cmd
		for _, match := range workerTaskModel.Matches {
			cmd := exec.Command("/bin/bash", "-c",
				fmt.Sprintf(`echo %s | su - %s -c "/usr/bin/bash %s/script/svn2git/svn2git.sh %s %s %s %s %s %d %t %s %s %s"`,
					workerTaskModel.CcPassword, workerTaskModel.CcUser, cwd, workerTaskModel.Pvob, workerTaskModel.Component,
					match.Stream, gitUrl, match.Branch, workerTaskModel.TaskId,
					workerTaskModel.IncludeEmpty, workerTaskModel.GitUser, workerTaskModel.GitEmail, workerTaskModel.Keep))
			cmds = append(cmds, cmd)
		}
		go infoServerTaskCompleted(&workerTaskModel, serverFlag, cmds)
		return http.StatusOK
	} else {
		log.Error(err)
		return http.StatusInternalServerError
	}
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
