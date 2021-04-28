package restapi

import (
	"ctgb/database"
	"ctgb/models"
	"ctgb/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	"net/http"
)

func CreateUser(params operations.CreateUserParams) middleware.Responder {
	var id int
	row := database.DB.QueryRow("SELECT id FROM user WHERE username=?", params.UserInfo.Username)
	err := row.Scan(&id)
	if err != nil || id != 0 {
		return middleware.Error(http.StatusInternalServerError, "User Already Exist")
	}
	sqlStr := "INSERT INTO user (username,password) VALUES (?,?)"
	ret := database.DB.MustExec(sqlStr, params.UserInfo.Username, params.UserInfo.Password)
	_, err = ret.RowsAffected()
	if err != nil {
		return middleware.Error(http.StatusInternalServerError, "Sql Error")
	} else {
		return middleware.Error(http.StatusCreated, "User Create Success")
	}
}

func ListUsers() middleware.Responder {
	rows, err := database.DB.Query("SELECT username FROM user")
	if err != nil {
		return middleware.Error(http.StatusInternalServerError, "Sql Error")
	}
	defer rows.Close()
	var users []*models.UserModel
	for rows.Next() {
		tmp := &models.UserModel{}
		if err := rows.Scan(&tmp.Username); err != nil {
			return middleware.Error(http.StatusInternalServerError, "Sql Error")
		}
		users = append(users, tmp)
	}
	return middleware.Error(http.StatusCreated, users)
}

func Login(params operations.LoginParams) middleware.Responder {
	var passwordInDB string
	row := database.DB.QueryRow("SELECT password FROM user WHERE username=?", params.UserInfo.Username)
	err := row.Scan(&passwordInDB)
	if err != nil {
		return middleware.Error(http.StatusInternalServerError, "")
	}
	if passwordInDB == params.UserInfo.Password {
		return middleware.Error(http.StatusCreated, "Login Success")
	} else {
		return middleware.Error(http.StatusInternalServerError, "User Does Not Exist Or Wrong Password")
	}
}
