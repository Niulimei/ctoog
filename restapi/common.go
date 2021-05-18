package restapi

import (
	"ctgb/utils"
	"net/http"
)

func CheckPermission(token string) bool {
	username, valid := utils.Verify(token)
	if !valid {
		utils.RecordLog(utils.Info, utils.Auth, "", "user "+username+" Unauthorized", http.StatusUnauthorized)
		return false
	}
	userInfo := getUserInfo(username)
	if userInfo.RoleID != int64(AdminRole) {
		utils.RecordLog(utils.Info, utils.Auth, "", "user "+username+" Forbidden", http.StatusForbidden)
		return false
	}
	return true
}
