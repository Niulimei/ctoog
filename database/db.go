package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var DB *sqlx.DB

var schema = `
CREATE TABLE user (
    id integer PRIMARY KEY autoincrement,
    username varchar (128),
    password varchar (32)
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
    last_completed_date_time varchar (64),
    creator varchar(128),
    worker_id integer
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

CREATE TABLE worker (
    id integer PRIMARY KEY autoincrement,
    worker_url varchar (256),
    status varchar (16),
    task_count integer
);
`

type TaskModel struct {
	// id
	Id int64 `db:"id"`

	// cc password
	CcPassword string `db:"cc_password"`

	// cc user
	CcUser string `json:"db:"cc_user"`

	// component
	Component string

	// git password
	GitPassword string `json:"db:"git_password"`

	// git URL
	GitURL string `db:"git_url"`

	// git user
	GitUser string `db:"git_user"`

	// git repo
	GitRepo string `db:"git_repo"`

	// status
	Status string

	//last completed date time
	LastCompletedDateTime string `db:"last_completed_date_time"`

	//worker id
	WorkerId int64 `db:"worker_id"`

	// pvob
	Pvob string
}

type WorkerModel struct {
	Id        int64
	WorkerUrl string `db:"worker_url"`
	Status    string
	TaskCount int64 `db:"task_count"`
}

type TaskLog struct {
	LogId     string `db:"log_id"`
	TaskId    int    `db:"task_id"`
	Status    string
	StartTime string `db:"start_time"`
	EndTime   string `db:"end_time"`
	Duration  int
}

func init() {
	var err error
	var isInitAlready = true
	_, err = os.Stat("translator.db") //os.Stat获取文件信息
	if err != nil {
		if os.IsNotExist(err) {
			isInitAlready = false
		}
	}
	DB, err = sqlx.Connect("sqlite3", "file:translator.db?cache=private&mode=rwc")
	if err != nil {
		log.Fatalln(err)
	}
	err = DB.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	if !isInitAlready {
		DB.MustExec(schema)
	}
}
