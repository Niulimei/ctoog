package restapi

import (
	"ctgb/database"
	"ctgb/restapi/operations"
	"fmt"
	"os/exec"
	"strings"

	"github.com/go-openapi/runtime/middleware"
)

var SvnUserName map[string][]string

func ProcessSvnUserName(svnUrl, svnUser, svnPassword string, taskId int64) {
	var names []string
	var checkedNames []string
	svn := fmt.Sprintf(`svn log --quiet --non-interactive --username "%s" --password "%s" "%s"`, svnUser, svnPassword, svnUrl)
	svnCmd := fmt.Sprintf(`%s | grep -a -E 'r[0-9]+ \| .+ \|' | cut -d '|' -f2 | sed 's/ //g' | sort | uniq`, svn)
	cmd := exec.Command("/bin/bash", "-c", svnCmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	result := string(out)
	//log.Debug("cmd", cmd.String(), "result:", result)
	names = strings.Split(result, "\n")
	for _, name := range names {
		if len(name) > 0 {
			checkedNames = append(checkedNames, name)
		}

	}
	if len(checkedNames) > 0 {
		SvnUserName[svnUrl + svnUser + svnPassword] = checkedNames
		if taskId != 0 {
			database.DB.Exec("UPDATE task SET status = 'init' WHERE id = ? and status = 'pending'", taskId)
		}
	}
}

func GetALlSvnName(svnUrl, svnUser, svnPassword string) []string {
	var checkedNames []string
	key := svnUrl + svnUser + svnPassword
	go ProcessSvnUserName(svnUrl, svnUser, svnPassword, 0)
	if storeNames, ok := SvnUserName[key]; ok {
		return storeNames
	}
	return checkedNames
}

func ListSvnUsernameHandler(params operations.ListSvnUsernameParams) middleware.Responder {
	return operations.NewListSvnUsernameOK().WithPayload(GetALlSvnName(params.SvnURL, params.SvnUser, params.SvnPassword))
}
func UpdateSvnNamePairHandler(params operations.UpdateSvnNamePairParams) middleware.Responder {
	return nil
}
