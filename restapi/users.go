package restapi

import (
	"ctgb/database"
	"ctgb/models"
	"ctgb/restapi/operations"
	"ctgb/utils"
	"github.com/go-openapi/runtime/middleware"
	"net/http"
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
	sqlStr := "INSERT INTO user (username,password) VALUES (?,?)"
	ret := database.DB.MustExec(sqlStr, params.UserInfo.Username, params.UserInfo.Password)
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

func checkPermission(token string) middleware.Responder {
	username, valid := utils.Verify(token)
	if !valid {
		return operations.NewListUserInternalServerError().WithPayload(&models.ErrorModel{
			Code:    http.StatusUnauthorized,
			Message: "",
		})
	}
	if username != "admin" {
		return operations.NewListUserInternalServerError().WithPayload(&models.ErrorModel{
			Code:    http.StatusForbidden,
			Message: "",
		})
	}
	return nil
}

func ListUsersHandler(param operations.ListUserParams) middleware.Responder {
	checkRet := checkPermission(param.Authorization)
	if checkRet != nil {
		return checkRet
	}
	rows, err := database.DB.Query("SELECT count(1) over() AS total_rows,username FROM user ORDER BY id LIMIT ? OFFSET ?", param.Limit, param.Offset)
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
		if err := rows.Scan(&count, &tmp.Username); err != nil {
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
			Code:    http.StatusInternalServerError,
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
			Code:    http.StatusInternalServerError,
			Message: "Wrong Password",
		})
	}
}
