package restapi

import (
	"ctgb/database"
	"ctgb/models"
	"ctgb/restapi/operations"
	"ctgb/utils"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

type ROLE int

const (
	AdminRole  ROLE = 1
	NormalRole ROLE = 2
)

func CreateUserHandler(params operations.CreateUserParams) middleware.Responder {
	if !CheckPermission(params.HTTPRequest) {
		return operations.NewCreateUserInternalServerError().WithPayload(&models.ErrorModel{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		})
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
		utils.RecordLog(utils.Info, utils.AddUser, "", "User "+params.UserInfo.Username+" create success", 0)
		return operations.NewCreateUserCreated().WithPayload(&models.OK{Message: "User " + params.UserInfo.Username + " create success"})
	}
}

func RegisterUserHandler(params operations.RegisterUserParams) middleware.Responder {
	var id int
	row := database.DB.QueryRow("SELECT id FROM user WHERE username=?", params.UserRegisterInfo.Username)
	err := row.Scan(&id)
	if err == nil || id != 0 {
		return operations.NewRegisterUserInternalServerError().WithPayload(&models.ErrorModel{
			Code:    http.StatusInternalServerError,
			Message: "Username Already Exist",
		})
	}
	sqlStr := "INSERT INTO user (bussinessgroup,team,nickname,username,password,role_id) VALUES (?,?,?,?,?,?)"
	ret := database.DB.MustExec(sqlStr, params.UserRegisterInfo.Bussinessgroup, params.UserRegisterInfo.Team,
		params.UserRegisterInfo.Nickname, params.UserRegisterInfo.Username, params.UserRegisterInfo.Password, NormalRole)
	_, err = ret.RowsAffected()
	if err != nil {
		return operations.NewRegisterUserInternalServerError().WithPayload(&models.ErrorModel{
			Code:    http.StatusInternalServerError,
			Message: "Sql Error",
		})
	} else {
		utils.RecordLog(utils.Info, utils.AddUser, "", "User "+params.UserRegisterInfo.Username+" register success", 0)
		return operations.NewRegisterUserCreated().WithPayload(&models.OK{Message: "User " + params.UserRegisterInfo.Username + " register success"})
	}
}

func getUserInfo(username string) *models.UserInfoModel {
	user := &models.UserInfoModel{}
	err := database.DB.Get(user, "SELECT * FROM user WHERE username=?", username)
	if err != nil {
		return user
	}
	return user
}

func GetUserHandler(param operations.GetUserParams) middleware.Responder {
	username := param.HTTPRequest.Header.Get("username")
	userInfo := getUserInfo(username)
	return operations.NewGetUserOK().WithPayload(&models.UserInfoModel{Username: username, RoleID: userInfo.RoleID})
}

func ListUsersHandler(param operations.ListUserParams) middleware.Responder {
	if !CheckPermission(param.HTTPRequest) {
		return operations.NewListUserInternalServerError().WithPayload(&models.ErrorModel{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}
	var users []*models.UserInfoModel
	err := database.DB.Select(&users, "SELECT * FROM user ORDER BY id LIMIT ? OFFSET ?", param.Limit, param.Offset)
	if err != nil {
		return operations.NewListUserInternalServerError().WithPayload(&models.ErrorModel{
			Code:    http.StatusInternalServerError,
			Message: "Sql Error",
		})
	}
	for i := range users {
		users[i].Password = ""
	}
	var count int64
	database.DB.Get(&count, "select count(id) from user")

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
		utils.RecordLog(utils.Error, utils.Login, "", "user "+params.UserInfo.Username+" does not exist.", http.StatusNotFound)
		return operations.NewLoginInternalServerError().WithPayload(&models.ErrorModel{
			Code:    http.StatusNotFound,
			Message: "User Does Not Exist",
		})
	}
	if passwordInDB == params.UserInfo.Password {
		token := utils.CreateJWT(params.UserInfo.Username)
		utils.RecordLog(utils.Info, utils.Login, "", "", 0)
		return operations.NewLoginCreated().WithPayload(&models.Authorization{
			Token: token,
		})
	} else {
		utils.RecordLog(utils.Error, utils.Login, "", "user "+params.UserInfo.Username+" password wrong.", http.StatusUnauthorized)
		return operations.NewLoginInternalServerError().WithPayload(&models.ErrorModel{
			Code:    http.StatusUnauthorized,
			Message: "Wrong Password",
		})
	}
}
