package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	dsn := "root:12345678@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB = db

	// 迁移 schema
	db.AutoMigrate(&History{})
}
