package restapi

import (
	"ctgb/database"
	"ctgb/restapi/operations"
	"database/sql"
	"encoding/json"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/go-openapi/runtime/middleware"
	log "github.com/sirupsen/logrus"
)

var pvobs []string

func init() {
	go func() {
		GetAllPvob()
		t := time.NewTicker(time.Second * 600)
		for {
			select {
			case <-t.C:
				GetAllPvob()
			}
		}
	}()

	go func() {
		for {
			updateStream()
			time.Sleep(time.Hour)
		}
	}()
}

func GetAllPvob() {
	pvobs = make([]string, 0)
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
	sort.Strings(pvobs)
	return
}

func RemoveDuplicatesAndEmpty(a []string) (ret []string){
	aLen := len(a)
	for i:=0; i < aLen; i++{
		if (i > 0 && a[i-1] == a[i]) || len(a[i])==0{
			continue
		}
		ret = append(ret, a[i])
	}
	return
}

func GetAllComponent(pvob string) []string {
	components := make([]string, 0)
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
	parseComponents := make([]string, 0)
	for _, component := range components {
		parseComponents = append(parseComponents, component)
		if strings.Count(component, "/") > 2 {
			tmp := strings.Split(component, "/")
			parseComponents = append(parseComponents, "/"+tmp[1]+"/"+tmp[2])
		}
	}
	sort.Strings(parseComponents)
	parseComponents = RemoveDuplicatesAndEmpty(parseComponents)
	return parseComponents
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
	streams := make([]string, 0)
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
	streams, err := getStreamFromDB(pvob, component)
	if err != nil && err != sql.ErrNoRows {
		return operations.NewListPvobComponentStreamInternalServerError()
	}
	if err == sql.ErrNoRows {
		streamsReal := GetAllStream(pvob, component)
		saveStreamToDB(pvob, component, streamsReal)
		streams = streamsReal
	}
	return operations.NewListPvobComponentStreamOK().WithPayload(streams)
}

func saveStreamToDB(pvob, component string, stream []string) error {
	streamStr, err := json.Marshal(stream)
	if err != nil {
		return err
	}
	_, err = database.DB.Exec("INSERT OR REPLACE INTO cc_repo (pvob, component, stream) VALUES (?,?,?)", pvob, component, streamStr)
	if err != nil {
		return err
	}
	return nil
}

func getStreamFromDB(pvob, component string) ([]string, error) {
	var total int
	var stream string
	err := database.DB.Get(&total, "SELECT count(1) from cc_repo where pvob=? and component=?", pvob, component)
	if err != nil && err != sql.ErrNoRows {
		log.Errorln(err.Error())
		return []string{}, err
	}
	if total != 1 {
		database.DB.Exec("DELETE FROM cc_repo where pvob=? and component=?", pvob, component)
		return []string{}, sql.ErrNoRows
	}
	err = database.DB.Get(&stream, "SELECT stream from cc_repo where pvob=? and component=?", pvob, component)
	if err != nil {
		log.Errorln(err.Error())
		return []string{}, err
	}

	var ret []string
	err = json.Unmarshal([]byte(stream), &ret)
	if err != nil {
		log.Errorln(err.Error())
		return []string{}, err
	} else {
		return ret, nil
	}
}

type PC struct {
	Pvob      string
	Component string
}

func updateStream() {
	var pc []PC
	err := database.DB.Select(&pc, "SELECT pvob, component FROM cc_repo")
	if err != nil {
		database.DB.Exec("DELETE FROM cc_repo")
		return
	}
	for _, v := range pc {
		err = saveStreamToDB(v.Pvob, v.Component, GetAllStream(v.Pvob, v.Component))
		if err != nil {
			database.DB.Exec("DELETE FROM cc_repo WHERE pvob=? AND component=?", v.Pvob, v.Component)
		}
	}
}
