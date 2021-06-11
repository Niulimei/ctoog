package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

type nodeInfo struct {
	IP         string `json:"IP"`
	User       string `json:"USER"`
	Password   string `json:"PASSWORD"`
	WorkDir    string `json:"WORKDIR"`
	ListenIP   string `json:"LISTEN_IP"`
	ListenPort int    `json:"LISTEN_PORT"`
}

type DeployConf struct {
	Server  *nodeInfo   `json:"SERVER"`
	Workers []*nodeInfo `json:"WORKERS"`
}

func main() {
	ret, err := ioutil.ReadFile("./conf.json")
	if err != nil {
		panic(err.Error())
	}
	var info = &DeployConf{}
	err = json.Unmarshal(ret, info)
	if err != nil {
		panic(err.Error())
	}
	if len(os.Args) != 2 {
		panic("Usage: ./deploy <version tar.gz file>")
	}
	versionFile := os.Args[1]
	//versionFile := "version_20210611154031.tar.gz"
	installServer := fmt.Sprintf(`bash -x ./deploy.sh "%s" "%s" "%s" "%s" "%s" "%s" "%s" "%d"`,
		info.Server.IP, info.Server.User, info.Server.Password, versionFile,
		info.Server.WorkDir, "translator-server", info.Server.ListenIP, info.Server.ListenPort)
	ret, err = exec.Command("/bin/bash", "-c", installServer).Output()
	if err != nil {
		fmt.Println("Install Server Error: ", err.Error())
	}
	fmt.Println("Install Server Result: ", string(ret))

	for _, worker := range info.Workers {
		installWorker := fmt.Sprintf(`bash -x ./deploy.sh "%s" "%s" "%s" "%s" "%s" "%s" "%s" "%d" "%s" "%d"`,
			worker.IP, worker.User, worker.Password, versionFile,
			worker.WorkDir, "translator-worker", info.Server.ListenIP, info.Server.ListenPort,
			worker.ListenIP, worker.ListenPort)
		ret, err = exec.Command("/bin/bash", "-c", installWorker).Output()
		if err != nil {
			fmt.Println("Install Worker Error: ", err.Error())
		}
		fmt.Println("Install Worker Result: ", string(ret))
	}
	fmt.Println("All Done!")
}
