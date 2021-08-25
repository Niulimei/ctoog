package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

//Global Var
var DB *gorm.DB

type History struct {
	gorm.Model
	Name        string
	Owner       string
	CreateTime  string
	HistoryId   string
	HistoryType string
	Description string
	GitName     string
}

func init() {
	var mysqlHost, mysqlPort, mysqlDatabase, mysqlUsername, mysqlPassword string
	mysqlHost, _ = os.LookupEnv("MYSQL_HOST")
	mysqlPort, _ = os.LookupEnv("MYSQL_PORT")
	mysqlDatabase, _ = os.LookupEnv("MYSQL_DATABASE")
	mysqlUsername, _ = os.LookupEnv("MYSQL_USERNAME")
	mysqlPassword, _ = os.LookupEnv("MYSQL_PASSWORD")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlUsername, mysqlPassword, mysqlHost, mysqlPort, mysqlDatabase)
	dsn = "root:12345678@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	fmt.Println("database init:", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB = db

	// 迁移 schema
	db.AutoMigrate(&History{})
}
