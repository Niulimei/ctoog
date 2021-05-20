package restapi

import (
	"ctgb/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func CheckPermission(r *http.Request) bool {
	username := r.Header.Get("username")
	userInfo := getUserInfo(username)
	if userInfo.RoleID != int64(AdminRole) {
		utils.RecordLog(utils.Info, utils.Auth, "", "user "+username+" Forbidden", http.StatusForbidden)
		return false
	}
	return true
}

var exceptionURL = map[string]string{
	"/login":   "POST",
	"/workers": "POST",
	"/tasks/cmdout": "POST",
}

func IsExceptionURL(method, uri string) bool {
	for uriChild, m := range exceptionURL {
		if strings.Contains(uri, uriChild) && m == method {
			log.Debug("got auth exception url", uri)
			return true
		}
	}
	return false
}
