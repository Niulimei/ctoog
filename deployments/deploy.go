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
	var svn, cc, jianxin string
	if _, exist := os.LookupEnv("SVN_SUPPORT"); exist {
		svn = "1"
	}
	if _, exist := os.LookupEnv("CC_SUPPORT"); exist {
		cc = "1"
	}
	if _, exist := os.LookupEnv("JIANXIN_SUPPORT"); exist {
		jianxin = "1"
	}
	installServer := fmt.Sprintf(`bash -x ./deploy.sh "%s" "%s" "%s" "%s" "%s" "%s" "%s" "%d" "%s" "%d" "%s" "%s" "%s"`,
		info.Server.IP, info.Server.User, info.Server.Password, versionFile,
		info.Server.WorkDir, "translator-server", info.Server.ListenIP, info.Server.ListenPort,
		"", 0, svn, cc, jianxin)
	//fmt.Println(installServer)
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
