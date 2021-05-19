package main

import (
	"bufio"
	"bytes"
	"ctgb/utils"
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
	"github.com/sevlyar/go-daemon"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var (
	hostFlag   string
	portFlag   int
	serverFlag string
)
var stop = make(chan struct{})

func init() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	// Only log the warning severity or above.
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	lvl = strings.ToLower(lvl)
	if ok {
		if lvl == "debug" {
			log.SetLevel(log.DebugLevel)
		} else if lvl == "info" {
			log.SetLevel(log.InfoLevel)
		} else {
			log.SetLevel(log.WarnLevel)
		}
	}
}

type commandOut struct {
	Logid   int64  `json:"log_id"`
	Content string `json:"content"`
}

type payload struct {
	Logid     string `json:"logID"`
	Status    string `json:"status"`
	Starttime string `json:"startTime"`
	Endtime   string `json:"endTime"`
	Duration  string `json:"duration"`
}

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
	req.Header.Set("Authorization", "1234567")

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.Body == nil {
		// handle err
		time.Sleep(time.Second * 3)
		resp, err = http.DefaultClient.Do(req)
		if err != nil || resp.Body == nil {
			return
		}
	}
	log.Info("info server success")
	defer resp.Body.Close()
}

func pingServer(host string, port int) {
	defer func() {
		if ret := recover(); ret != nil {
			fmt.Printf("Recover From Panic. %v\n", ret)
		}
	}()
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

type MatchInfo struct {
	Branch string
	Stream string
}

type Task struct {
	TaskId       int64
	TaskLogId    int64
	CcPassword   string
	CcUser       string
	Component    string
	GitPassword  string
	GitURL       string
	GitUser      string
	GitEmail     string
	Pvob         string
	IncludeEmpty bool
	Matches      []MatchInfo
	Keep         string
}

type TaskDelInfo struct {
	TaskId     int64  `json:"task_id"`
	CcPassword string `json:"cc_password"`
	CcUser     string `json:"cc_user"`
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
		cmdStr := fmt.Sprintf(`echo %s | su - %s -c "/usr/bin/bash %s/cleanCache.sh %d"`,
			workerTaskModel.CcPassword, workerTaskModel.CcUser, cwd, workerTaskModel.TaskId)
		fmt.Println(cmdStr)
		cmd := exec.Command("/bin/bash", "-c", cmdStr)
		out, err := cmd.Output()
		if err != nil {
			log.Error(err)
			w.WriteHeader(500)
			w.Write(out)
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
		var cmds []*exec.Cmd
		for _, match := range workerTaskModel.Matches {
			cmd := exec.Command("/bin/bash", "-c",
				fmt.Sprintf(`echo %s | su - %s -c "/usr/bin/bash %s/cc2git.sh %s %s %s %s %s %d %t %s %s %s"`,
					workerTaskModel.CcPassword, workerTaskModel.CcUser, cwd, workerTaskModel.Pvob, workerTaskModel.Component,
					match.Stream, gitUrl, match.Branch, workerTaskModel.TaskId,
					workerTaskModel.IncludeEmpty, workerTaskModel.GitUser, workerTaskModel.GitEmail, workerTaskModel.Keep))
			cmds = append(cmds, cmd)
		}
		go infoServerTaskCompleted(&workerTaskModel, serverFlag, cmds)
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
	usage := "Usage: ./translator-worker start | stop | status | restart"
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "start":
			fmt.Println("translator-worker is starting.")
			break
		case "stop", "restart":
			pid, _ := ioutil.ReadFile("./translator-worker.pid")
			if string(pid) == "" {
				fmt.Println("translator-worker is not running.")
			} else if _, err := utils.Exec("ps " + string(pid)); err != nil {
				fmt.Println("translator-worker is not running.")
			} else {
				utils.Exec("kill " + string(pid))
				fmt.Println("translator-worker has been stopped.")
			}
			if command == "stop" {
				return
			} else {
				fmt.Println("translator-worker is starting.")
				break
			}
		case "status":
			pid, _ := ioutil.ReadFile("./translator-worker.pid")
			if string(pid) == "" {
				fmt.Println("translator-worker is not running.")
			} else if _, err := utils.Exec("ps " + string(pid)); err != nil {
				fmt.Println("translator-worker is not running.")
			} else {
				fmt.Println("translator-worker is running.")
			}
			return
		default:
			fmt.Println(usage)
			return
		}
	}
	cntxt := &daemon.Context{
		PidFileName: "translator-worker.pid",
		PidFilePerm: 0644,
		LogFileName: "translator-worker.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"translator-worker"},
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	log.Print("- - - - - - - - - - - - - - -")
	log.Print("daemon started")
	type conf struct {
		Host       string `yaml:"host"`
		Port       int    `yaml:"port"`
		ServerAddr string `yaml:"server_addr"`
	}
	tmp := &conf{}
	content, _ := ioutil.ReadFile("./translator-worker.yaml")
	yaml.Unmarshal(content, tmp)
	serverFlag = tmp.ServerAddr
	hostFlag = tmp.Host
	portFlag = tmp.Port
	go pingServer(hostFlag, portFlag)
	http.HandleFunc("/new_task", taskHandler) //	设置访问路由
	http.HandleFunc("/delete_task", deleteTaskHandler)
	log.Fatal(http.ListenAndServe(hostFlag+":"+strconv.Itoa(portFlag), nil))
}
