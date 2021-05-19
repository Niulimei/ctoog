package restapi

import (
	"ctgb/utils"
	"net/http"
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
