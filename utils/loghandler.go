package utils

import (
	"ctgb/database"
	"time"
)

const (
	Fatal LogLevel = iota + 1
	Error
	Warning
	Info
	Debug
)

const (
	GetTask    Action = "GET_TASK"
	AddTask    Action = "ADD_TASK"
	UpdateTask Action = "UPDATE_TASK"
	DeleteTask Action = "DELETE_TASK"
	GetUser    Action = "GET_USER"
	AddUser    Action = "ADD_USER"
	UpdateUser Action = "UPDATE_USER"
	DeleteUser Action = "DELETE_USER"
	Login      Action = "LOGIN"
)

type Action string

type ErrorCode int

type LogLevel int

type Log struct {
	Time     int64
	Level    LogLevel
	User     string
	Action   Action
	Position string
	Message  string
	ErrCode  ErrorCode
}

var LogChannel = make(chan *Log, 100)

func LogHandle() {
	sqlStr := "INSERT INTO log (time, level, user, action, position, message, errcode) VALUES(?,?,?,?,?,?,?)"
	for {
		select {
		case logContent := <-LogChannel:
			database.DB.MustExec(sqlStr, logContent.Time, logContent.Level, logContent.User,
				logContent.Action, logContent.Position, logContent.Message, logContent.ErrCode)
		}
	}
}

func RecordLog(level LogLevel, action Action, position, message string, errCode ErrorCode) {
	in := &Log{
		Time:     time.Now().Unix(),
		Level:    level,
		User:     "admin",
		Action:   action,
		Position: position,
		Message:  message,
		ErrCode:  errCode,
	}
	LogChannel <- in
}
