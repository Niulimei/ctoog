package restapi

import (
	"ctgb/models"
	"ctgb/restapi/operations"
	"ctgb/utils"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

func CheckPermission(token string) middleware.Responder {
	username, valid := utils.Verify(token)
	if !valid {
		utils.RecordLog(utils.Info, utils.Login, "", "user "+username+" Unauthorized", http.StatusUnauthorized)
		return operations.NewListUserInternalServerError().WithPayload(&models.ErrorModel{
			Code:    http.StatusUnauthorized,
			Message: "",
		})
	}
	userInfo := getUserInfo(username)
	if userInfo.RoleID != int64(AdminRole) {
		utils.RecordLog(utils.Info, utils.Login, "", "user "+username+" Forbidden", http.StatusForbidden)
		return operations.NewListUserInternalServerError().WithPayload(&models.ErrorModel{
			Code:    http.StatusForbidden,
			Message: "",
		})
	}
	return nil
}
