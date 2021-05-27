package restapi

import (
	"ctgb/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	"os/exec"
	"strings"
)

func GetALlSvnName(svnUrl, svnUser, svnPassword string) []string {
	tmp := strings.SplitN(svnUrl, "//", 2)
	svnUrl = tmp[0] + "//" + svnUser + ":" + svnPassword + "@" + tmp[1]
	cmd := exec.Command("svn", "log", "--quiet", svnUrl)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil
	}
	result := string(out)
	//log.Debug("cmd", cmd.String(), "result:", result)
	names := strings.Split(result, "\n")
	return names
}

func ListSvnUsernameHandler(params operations.ListSvnUsernameParams) middleware.Responder {
	return operations.NewListSvnUsernameOK().WithPayload(GetALlSvnName(params.SvnURL, params.SvnUser, params.SvnPassword))
}
func UpdateSvnNamePairHandler(params operations.UpdateSvnNamePairParams) middleware.Responder {
	return nil
}
