package restapi

import (
	"ctgb/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strings"
	"time"
)

var pvobs []string

func init() {
	go func() {
		GetAllPvob()
		t := time.NewTimer(time.Second * 600)
		for {
			select {
			case <-t.C:
				GetAllPvob()
				t.Reset(time.Second * 600)
			}
		}
	}()
}

func GetAllPvob() {
	pvobs = make([]string, 0, 10)
	cmd := exec.Command("cleartool", "lsvob", "-l")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	result := string(out)
	//log.Debug("cmd", cmd.String(), "result:", result)
	infos := strings.Split(result, "\n\n")
	for _, info := range infos {
		lines := strings.Split(info, "\n")
		for i, l := range lines {
			log.Debug(i, l)
		}
		if len(lines) == 0 {
			continue
		}
		log.Debug(lines[0], lines[len(lines)-1], strings.HasPrefix(lines[0], "Tag: "), lines[len(lines)-1] == "Vob registry attributes: ucmvob")
		if strings.HasPrefix(lines[0], "Tag: ") && lines[len(lines)-1] == "Vob registry attributes: ucmvob" {
			pvob := lines[0][5:]
			pvobs = append(pvobs, pvob)
		}
	}
	return
}

func GetAllComponent(pvob string) []string {
	components := make([]string, 0, 10)
	args := `lscomp -fmt "%[root_dir]p\n" -invob ` + pvob
	cmd := exec.Command("cleartool", strings.Split(args, " ")...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil
	}
	result := string(out)
	result = strings.Replace(result, `"`, "", -1)
	log.Debug("cmd", cmd.String(), "result:", result)
	lines := strings.Split(result, "\n")
	for _, line := range lines {
		// cleartool 命令返回的信息里可能存在cleartool的提示或者警告信息，不是我们期望的内容，以 cleartoo: 开头，应该跳过
		if len(line) > 0 && strings.Index(line, "cleartool: ") == -1 {
			components = append(components, line)
		}
	}
	return components
}

func checkStreamComponent(pvob, component, stream string) bool {
	log.Debug("stream:", pvob, ":", component, ":", stream)
	args := `lsstream -fmt %[components]p ` + stream + "@" + pvob
	cmd := exec.Command("cleartool", strings.Split(args, " ")...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("cmd:", cmd.String())
		log.Error("stream list comp:", err)
		return false
	}
	result := string(out)
	log.Debug("cmd", cmd.String(), "result:", result)
	lines := strings.Split(result, "\n")
	tmp := strings.Split(component, "/")
	component = tmp[len(tmp)-1]
	for _, line := range lines {
		if line == component {
			return true
		}
		// cleartool 命令返回的信息里可能存在cleartool的提示或者警告信息，不是我们期望的内容，以 cleartoo: 开头，应该跳过
		log.Debug("line", line, strings.Index(line, "cleartool: "))
		if strings.Index(line, "cleartool: ") == -1 {
			ls := strings.Split(line, " ")
			log.Debug("ls", ls)
			for _, l := range ls {
				if l == component {
					return true
				}
			}
		}
	}
	return false
}

func GetAllStream(pvob, component string) []string {
	streams := make([]string, 0, 10)
	args := `lsstream -s -invob ` + pvob
	cmd := exec.Command("cleartool", strings.Split(args, " ")...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil
	}
	result := string(out)
	log.Debug("cmd", cmd.String(), "result:", result)
	lines := strings.Split(result, "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		// cleartool 命令返回的信息里可能存在cleartool的提示或者警告信息，不是我们期望的内容，以 cleartoo: 开头，应该跳过
		if strings.Index(line, "cleartool: ") != -1 {
			continue
		}
		if checkStreamComponent(pvob, component, line) {
			streams = append(streams, line)
		}
	}
	return streams
}

func ListPvobHandler(params operations.ListPvobParams) middleware.Responder {
	return operations.NewListPvobOK().WithPayload(pvobs)
}

func ListPvobComponentHandler(params operations.ListPvobComponentParams) middleware.Responder {
	pvob := params.ID
	return operations.NewListPvobComponentOK().WithPayload(GetAllComponent(pvob))
}

func ListPvobComponentStreamHandler(params operations.ListPvobComponentStreamParams) middleware.Responder {
	pvob := params.PvobID
	component := params.ComponentID
	return operations.NewListPvobComponentStreamOK().WithPayload(GetAllStream(pvob, component))
}
