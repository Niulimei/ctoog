package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/urlesc"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	l "log"
	"net/http"
	"os"
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

func init() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	if ok {
		if strings.ToLower(lvl) == "debug" {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.WarnLevel)
		}
	}
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
	log.Debug("start cmd:", cmd.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		end := time.Now()
		duration := end.Sub(start).Seconds()
		d := strconv.FormatInt(int64(duration), 10)
		data.Status = "failed"
		data.Duration = d
	}

	result := string(out)
	log.Debug(result)

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
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/api/workers", serverFlag), body)
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
	TaskId       int64
	CcPassword   string
	CcUser       string
	Component    string
	GitPassword  string
	GitURL       string
	GitUser      string
	GitEmail     string
	Pvob         string
	Stream       string
	Branch       string
	IncludeEmpty bool
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	workerTaskModel := Task{}
	if err := json.Unmarshal(body, &workerTaskModel); err == nil {
		workerTaskModel.GitUser = urlesc.QueryEscape(workerTaskModel.GitUser)
		workerTaskModel.GitPassword = urlesc.QueryEscape(workerTaskModel.GitPassword)
		log.Printf("%+v\n", workerTaskModel)
		gitUrl := workerTaskModel.GitURL
		if strings.HasPrefix(gitUrl, "http://") {
			gitUrl = strings.Replace(gitUrl, "http://", "", 1)
			gitUrl = "http://" + workerTaskModel.GitUser + ":" + workerTaskModel.GitPassword + "@" + gitUrl
		} else if strings.HasPrefix(gitUrl, "https://") {
			gitUrl = strings.Replace(gitUrl, "https://", "", 1)
			gitUrl = "https://" + workerTaskModel.GitUser + ":" + workerTaskModel.GitPassword + "@" + gitUrl
		}
		cwd, _ := os.Getwd()
		cmd := exec.Command("/bin/bash", "-c",
			fmt.Sprintf(`echo %s | su - %s -c "/usr/bin/bash %s/cc2git.sh %s %s %s %s %s %d %t %s %s"`,
				workerTaskModel.CcPassword, workerTaskModel.CcUser, cwd, workerTaskModel.Pvob, workerTaskModel.Component,
				workerTaskModel.Stream, gitUrl, workerTaskModel.Branch, workerTaskModel.TaskId,
				workerTaskModel.IncludeEmpty, workerTaskModel.GitUser, workerTaskModel.GitEmail))
		go infoServerTaskCompleted(&workerTaskModel, serverFlag, cmd)
	} else {
		log.Error(err)
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
	l.Fatal(http.ListenAndServe(hostFlag+":"+strconv.Itoa(portFlag), nil))
}
