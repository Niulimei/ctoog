package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path/filepath"
)

var DB *sqlx.DB

var schema = `
CREATE TABLE user (
    id integer PRIMARY KEY autoincrement,
    username varchar (50),
    password_hash varchar (32)
);

CREATE TABLE task (
    id integer PRIMARY KEY autoincrement,
    pvob varchar (256),
    component varchar (256),
    cc_user varchar (256),
    cc_password varchar (256),
    git_url varchar (256),
    git_user varchar (256),
    git_password varchar (256),
    git_repo varchar (256),
    status varchar (16),
    last_completed_date_time varchar (64)
);

CREATE TABLE match_info (
    id integer PRIMARY KEY autoincrement,
    task_id integer,
    stream varchar (256),
    git_branch varchar (256)
);

CREATE TABLE task_log (
    log_id integer PRIMARY KEY autoincrement,
    task_id integer,
    status varchar (16),
    start_time varchar (64),
    end_time varchar (64),
    duration integer
);
`

func init() {
	var err error
	var isInitAlready bool
	dbDirPath, _ := os.Getwd()
	dbFilePath := filepath.Join(dbDirPath, "translator.db")
	fmt.Println(dbFilePath)
	_, err = os.Stat(dbFilePath) //os.Stat获取文件信息
	if err == nil {
		isInitAlready = true
	}
	DB, err = sqlx.Connect("sqlite3", fmt.Sprintf("file:%s?cache=private&mode=rwc", dbFilePath))
	if err != nil {
		log.Fatalln(err)
	}
	if !isInitAlready {
		DB.MustExec(schema)
	}
}
