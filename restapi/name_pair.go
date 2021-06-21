package restapi

import (
	"ctgb/restapi/operations"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"os/exec"
	"strings"
)

func GetALlSvnName(svnUrl, svnUser, svnPassword string) []string {
	tmp := strings.SplitN(svnUrl, "//", 2)
	var names []string
	var checkedNames []string
	if len(tmp) == 2 {
		svnUrl = tmp[0] + "//" + svnUser + ":" + svnPassword + "@" + tmp[1]
		svnCmd := fmt.Sprintf(`svn log --quiet %s | grep -a -E 'r[0-9]+ \| .+ \|' | cut -d '|' -f2 | sed 's/ //g' | sort | uniq`, svnUrl)
		cmd := exec.Command("/bin/bash", "-c", svnCmd)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return nil
		}
		result := string(out)
		//log.Debug("cmd", cmd.String(), "result:", result)
		names = strings.Split(result, "\n")
		for _, name := range names {
			if len(name) > 0 {
				checkedNames = append(checkedNames, name)
			}
		}
	}
	return checkedNames
}

func ListSvnUsernameHandler(params operations.ListSvnUsernameParams) middleware.Responder {
	return operations.NewListSvnUsernameOK().WithPayload(GetALlSvnName(params.SvnURL, params.SvnUser, params.SvnPassword))
}
func UpdateSvnNamePairHandler(params operations.UpdateSvnNamePairParams) middleware.Responder {
	return nil
}
