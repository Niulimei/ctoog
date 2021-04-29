package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	portFlag int
)

func init() {
	flag.IntVar(&portFlag, "port", 8080, "service listens on this port")
}

func infoServerTaskCompleted(task *Task) {
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
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		// handle err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("PUT",
		fmt.Sprintf("http://127.0.0.1:8993/api/tasks/%d?start=false", task.TaskId), body)
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

func pingServer() {

	type Payload struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}

	data := Payload{
		Host: "127.0.0.1",
		Port: 8995,
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
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	workerTaskModel := Task{}
	if err := json.Unmarshal(body, &workerTaskModel); err == nil {
		//TODO task process
	} else {
		fmt.Println(err)
		w.WriteHeader(500)
		w.Write([]byte("bye"))
		return
	}
	w.WriteHeader(201)
	w.Write([]byte("bye"))
	go infoServerTaskCompleted(&workerTaskModel)
}

func main() {
	flag.Parse()
	go pingServer()
	http.HandleFunc("/new_task", taskHandler) //	设置访问路由
	log.Fatal(http.ListenAndServe(strconv.Itoa(portFlag), nil))
}
