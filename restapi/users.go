package restapi

import (
	"ctgb/database"
	"ctgb/models"
	"ctgb/restapi/operations"
	"ctgb/utils"
	"github.com/go-openapi/runtime/middleware"
	"net/http"
)

type ROLE int

const (
	AdminRole  ROLE = 1
	NormalRole ROLE = 2
)

func CreateUserHandler(params operations.CreateUserParams) middleware.Responder {
	checkRet := checkPermission(params.Authorization)
	if checkRet != nil {
		return checkRet
	}
	var id int
	row := database.DB.QueryRow("SELECT id FROM user WHERE username=?", params.UserInfo.Username)
	err := row.Scan(&id)
	if err == nil || id != 0 {
		return operations.NewCreateUserInternalServerError().WithPayload(&models.ErrorModel{
			Code:    http.StatusInternalServerError,
			Message: "User Already Exist",
		})
	}
	sqlStr := "INSERT INTO user (username,password,role_id) VALUES (?,?,?)"
	if params.UserInfo.RoleID == 0 {
		params.UserInfo.RoleID = int64(NormalRole)
	}
	ret := database.DB.MustExec(sqlStr, params.UserInfo.Username, params.UserInfo.Password, params.UserInfo.RoleID)
	_, err = ret.RowsAffected()
	if err != nil {
		return operations.NewCreateUserInternalServerError().WithPayload(&models.ErrorModel{
			Code:    http.StatusInternalServerError,
			Message: "Sql Error",
		})
	} else {
		return operations.NewCreateUserCreated().WithPayload(&models.OK{Message: "User Create Success"})
	}
}

func getUserInfo(username string) *models.UserInfoModel {
	user := &models.UserInfoModel{}
	row := database.DB.QueryRow("SELECT username,role_id FROM user WHERE username=?", username)
	err := row.Scan(&user.Username, &user.RoleID)
	if err != nil {
		return user
	}
	return user
}

func checkPermission(token string) middleware.Responder {
	username, valid := utils.Verify(token)
	if !valid {
		return operations.NewListUserInternalServerError().WithPayload(&models.ErrorModel{
			Code:    http.StatusUnauthorized,
			Message: "",
		})
	}
	userInfo := getUserInfo(username)
	if userInfo.RoleID != int64(AdminRole) {
		return operations.NewListUserInternalServerError().WithPayload(&models.ErrorModel{
			Code:    http.StatusForbidden,
			Message: "",
		})
	}
	return nil
}

func GetUserHandler(param operations.GetUserParams) middleware.Responder {
	username, valid := utils.Verify(param.Authorization)
	if !valid {
		return operations.NewGetUserInternalServerError().WithPayload(&models.ErrorModel{
			Code:    http.StatusUnauthorized,
			Message: "",
		})
	}
	userInfo := getUserInfo(username)
	return operations.NewGetUserOK().WithPayload(&models.UserInfoModel{Username: username, RoleID: userInfo.RoleID})
}

func ListUsersHandler(param operations.ListUserParams) middleware.Responder {
	checkRet := checkPermission(param.Authorization)
	if checkRet != nil {
		return checkRet
	}
	rows, err := database.DB.Query("SELECT count(1) over() AS total_rows,username,role_id FROM user ORDER BY id LIMIT ? OFFSET ?", param.Limit, param.Offset)
	if err != nil {
		return operations.NewListUserInternalServerError().WithPayload(&models.ErrorModel{
			Code:    http.StatusInternalServerError,
			Message: "Sql Error",
		})
	}
	defer rows.Close()
	var users []*models.UserInfoModel
	var count int64
	for rows.Next() {
		tmp := &models.UserInfoModel{}
		if err := rows.Scan(&count, &tmp.Username, &tmp.RoleID); err != nil {
			return operations.NewListUserInternalServerError().WithPayload(&models.ErrorModel{
				Code:    http.StatusInternalServerError,
				Message: "Sql Error",
			})
		}
		users = append(users, tmp)
	}
	return operations.NewListUserOK().WithPayload(&models.UserPageInfoModel{
		Count:    count,
		Limit:    param.Limit,
		Offset:   param.Offset,
		UserInfo: users,
	})
}

func LoginHandler(params operations.LoginParams) middleware.Responder {
	var passwordInDB string
	row := database.DB.QueryRow("SELECT password FROM user WHERE username=?", params.UserInfo.Username)
	err := row.Scan(&passwordInDB)
	if err != nil {
		return operations.NewLoginInternalServerError().WithPayload(&models.ErrorModel{
			Code:    http.StatusNotFound,
			Message: "User Does Not Exist",
		})
	}
	if passwordInDB == params.UserInfo.Password {
		token := utils.CreateJWT(params.UserInfo.Username)
		return operations.NewLoginCreated().WithPayload(&models.Authorization{
			Token: token,
		})
	} else {
		return operations.NewLoginInternalServerError().WithPayload(&models.ErrorModel{
			Code:    http.StatusUnauthorized,
			Message: "Wrong Password",
		})
	}
}
