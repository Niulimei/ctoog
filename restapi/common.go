package restapi

import (
	"ctgb/utils"
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
	"/login": "POST",
}

func IsExceptionURL(method, uri string) bool {
	for uriChild, m := range exceptionURL {
		if strings.Contains(uri, uriChild) && m == method {
			return true
		}
	}
	return false
}
