package main

import (
	"ctgb/restapi"
	"ctgb/utils"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/sevlyar/go-daemon"
	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"
)

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
				os.RemoveAll("./translator-worker.pid")
				fmt.Println("translator-worker has been stopped.")
			}
			if command == "stop" {
				return
			} else {
				fmt.Println("translator-worker is starting.")
				time.Sleep(time.Millisecond * 500)
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
	logFile := "log/translator-worker.log"
	os.MkdirAll(filepath.Dir(logFile), 0666)
	cntxt := &daemon.Context{
		PidFileName: "translator-worker.pid",
		PidFilePerm: 0644,
		LogFileName: logFile,
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

	tmp := &restapi.Conf{}
	content, _ := ioutil.ReadFile("./translator-worker.yaml")
	yaml.Unmarshal(content, tmp)
	restapi.ServerFlag = tmp.ServerAddr
	go restapi.DumpLogFile(logFile)
	go restapi.CleanOldTmpCmdOutFile()
	go restapi.PingServer(tmp.Host, tmp.Port)
	http.HandleFunc("/new_task", restapi.WorkerTaskHandler) //	设置访问路由
	http.HandleFunc("/delete_task", restapi.DeleteWorkerTaskCacheHandler)
	http.HandleFunc("/command_out", restapi.GetCommandOut)
	http.HandleFunc("/check_info", restapi.CheckInfoHandler)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(tmp.Port), nil))
}
