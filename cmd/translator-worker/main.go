package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	hostFlag   string
	portFlag   int
	serverFlag string
)

func init() {
	flag.StringVar(&hostFlag, "host", "127.0.0.1", "service listens on this IP")
	flag.IntVar(&portFlag, "port", 8080, "service listens on this port")
	flag.StringVar(&serverFlag, "serverAddr", "127.0.0.1", "translator server listens on this IP:port")
}

func infoServerTaskCompleted(task *Task, server string, cmd *exec.Cmd) {

	type Payload struct {
		Logid     string `json:"logID"`
		Status    string `json:"status"`
		Starttime string `json:"startTime"`
		Endtime   string `json:"endTime"`
		Duration  string `json:"duration"`
	}

	data := Payload{
		Status:   "completed",
		Endtime:  time.Now().Format("2006-01-02 15:04:05"),
		Duration: "10",
	}
	start := time.Now()
	if err := cmd.Start(); err != nil {
		end := time.Now()
		duration := end.Sub(start).Seconds()
		d := strconv.FormatInt(int64(duration), 10)
		data.Status = "failed"
		data.Duration = d
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		// handle err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("PUT",
		fmt.Sprintf("http://%s/api/tasks/%d?start=false", server, task.TaskId), body)
	if err != nil {
		// handle err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "1234567")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
		time.Sleep(time.Second * 3)
		http.DefaultClient.Do(req)
	}
	defer resp.Body.Close()
}

func pingServer(host string, port int) {

	type Payload struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}

	data := Payload{
		Host: host,
		Port: port,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		// handle err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "http://127.0.0.1:8993/api/workers", body)
	if err != nil {
		// handle err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		// handle err
	}
	defer resp.Body.Close()
}

type Task struct {
	TaskId      int64
	CcPassword  string
	CcUser      string
	Component   string
	GitPassword string
	GitURL      string
	GitUser     string
	Pvob        string
	Stream      string
	Branch      string
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	workerTaskModel := Task{}
	gitUrl := workerTaskModel.GitURL
	if strings.HasPrefix(gitUrl, "http://") {
		gitUrl = strings.Replace(gitUrl, "http://", "", 1)
		gitUrl = "http://" + workerTaskModel.GitUser + ":" + workerTaskModel.GitPassword + "@" + gitUrl
	} else if strings.HasPrefix(gitUrl, "https://") {
		gitUrl = strings.Replace(gitUrl, "https://", "", 1)
		gitUrl = "https://" + workerTaskModel.GitUser + ":" + workerTaskModel.GitPassword + "@" + gitUrl
	}
	if err := json.Unmarshal(body, &workerTaskModel); err == nil {
		cmd := exec.Command("/bin/bash", "-c",
			fmt.Sprintf(`echo %s | sudo -S su - %s -c "/usr/bin/bash cc2git.sh" %s %s %s %s %s %d`,
				workerTaskModel.CcPassword, workerTaskModel.CcUser, workerTaskModel.Pvob, workerTaskModel.Component,
				workerTaskModel.Stream,
				gitUrl, workerTaskModel.Branch, workerTaskModel.TaskId))
		go infoServerTaskCompleted(&workerTaskModel, serverFlag, cmd)
	} else {
		fmt.Println(err)
		w.WriteHeader(500)
		w.Write([]byte("bye"))
		return
	}
	w.WriteHeader(201)
	w.Write([]byte("bye"))
}

func main() {
	flag.Parse()
	go pingServer(hostFlag, portFlag)
	http.HandleFunc("/new_task", taskHandler) //	设置访问路由
	log.Fatal(http.ListenAndServe(strconv.Itoa(portFlag), nil))
}
