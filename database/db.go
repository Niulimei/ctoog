package database

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

var DB *sqlx.DB

var schema = `
CREATE TABLE user (
    id integer PRIMARY KEY autoincrement,
    username varchar (128),
    password varchar (32),
    role_id integer
);

CREATE TABLE log (
    id integer PRIMARY KEY autoincrement,
    time integer,
	level varchar (256),
	user varchar (256),
	action varchar (256),
	position varchar (256),
	message varchar (256),
	errcode integer
);

CREATE TABLE task (
    id integer PRIMARY KEY autoincrement,
    cc_user varchar (256),
    cc_password varchar (256),
    git_url varchar (256),
    git_user varchar (256),
    git_password varchar (256),
    git_email varchar (256),
    status varchar (16),
    last_completed_date_time varchar (64),
    creator varchar(128),
    worker_id integer,
    model_type varchar(16) default 'clearcase'
);


CREATE TABLE cc_meta (
    id integer PRIMARY KEY autoincrement,
    pvob varchar (256),
    component varchar (256),
    dir varchar (256),
    keep varchar (256),
    task_id integer,
    include_empty boolean
);

CREATE TABLE svn_meta (
    id integer PRIMARY KEY autoincrement,
    svn_url varchar(512)
);

CREATE TABLE match_info (
    id integer PRIMARY KEY autoincrement,
    task_id integer,
    stream varchar (256),
    git_branch varchar (256)
);

CREATE TABLE svn_name_pair (
    id integer PRIMARY KEY autoincrement,
    task_id integer,
    svn_username varchar (256),
    git_username varchar (256),
    git_email varchar (256)
);

CREATE TABLE task_log (
    log_id integer PRIMARY KEY autoincrement,
    task_id integer,
    status varchar (16),
    start_time varchar (64),
    end_time varchar (64),
    duration integer
);

CREATE TABLE task_command_out (
    log_id integer PRIMARY KEY,
    content text
);

CREATE TABLE worker (
    id integer PRIMARY KEY autoincrement,
    worker_url varchar (256),
    status varchar (16),
    task_count integer,
    register_time varchar (64)
);

CREATE TABLE schedule (
    id integer PRIMARY KEY autoincrement,
    status varchar (16),
    schedule varchar (16),
    task_id integer,
    creator varchar (128)
);


INSERT INTO user (username,password,role_id) VALUES("admin", "b17eccdc6c06bd8e15928d583503adf9", 1);
`

type TaskModel struct {
	// id
	Id int64 `db:"id"`

	// cc password
	CcPassword string `db:"cc_password"`

	// cc user
	CcUser string `db:"cc_user"`

	// git password
	GitPassword string `db:"git_password"`

	// git URL
	GitURL string `db:"git_url"`

	// git user
	GitUser string `db:"git_user"`

	// status
	Status string

	//last completed date time
	LastCompletedDateTime string `db:"last_completed_date_time"`

	//worker id
	WorkerId int64 `db:"worker_id"`

	// git email
	GitEmail string `db:"git_email"`

	// model type
	ModelType string `db:"model_type"`
}

type CcMetaModel struct {
	// id
	Id int64

	// component
	Component string

	// dir
	Dir string

	// keep
	Keep string

	// pvob
	Pvob string

	// include empty dir
	IncludeEmpty bool `db:"include_empty"`
}

type SvnMetaModel struct {
	Id int64
	SvnURL string `db:"svn_url"`
}

type SvnNamePair struct {
	TaskId int64 `db:"task_id"`
	SvnUserName string `db:"svn_username"`
	GitUserName string `db:"git_username"`
	GitEmail    string `db:"git_email"`
}

type WorkerModel struct {
	Id           int64
	WorkerUrl    string `db:"worker_url"`
	Status       string
	TaskCount    int64  `db:"task_count"`
	RegisterTime string `db:"register_time"`
}

type TaskLog struct {
	LogId     string `db:"log_id"`
	TaskId    int    `db:"task_id"`
	Status    string
	StartTime string `db:"start_time"`
	EndTime   string `db:"end_time"`
	Duration  int
}

type MatchInfo struct {
	Id        string `db:"id"`
	TaskId    string `db:"task_id"`
	Stream    string `db:"stream"`
	GitBranch string `db:"git_branch"`
}

type ScheduleModel struct {
	Id       int64
	Status   string
	TaskId   int64 `db:"task_ik"`
	Schedule string
	Creator  string
}

func init() {
	initDB()
	upgradeDB()
	go func() {
		for {
			startTime := time.Now().Format("2006.01.02-15:04:05")
			cmd := exec.Command("cp", "translator.db", "backup/translator-"+startTime+".back")
			cmd.Run()
			time.Sleep(time.Minute * 10)
		}
	}()
	go func() {
		pattern := "backup/translator-*.back"
		for {
			paths, err := filepath.Glob(pattern)
			if err == nil {
				for _, path := range paths {
					d, err := time.ParseInLocation("backup/translator-2006.01.02-15:04:05.back", path, time.Local)
					if err == nil {
						duration := time.Now().Sub(d)
						if duration > time.Hour*24*15 {
							cmd := exec.Command("rm", path)
							err := cmd.Run()
							if err != nil {
								log.Error("delete backuo err:", err, " ", path)
							} else {
								log.Info("delete backup:", path)
							}
						}
					}
				}
			}
			time.Sleep(time.Hour)
		}
	}()
}

func initDB() {
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
		_, err = sqlx.LoadFile(DB, "sql/base.sql")
		if err != nil {
			log.Fatalln(err)
		}
	}
}

type sqlFiles struct {
	Files []string `json:"files"`
}

func upgradeDB() {
	var sf = &sqlFiles{}
	c, err := ioutil.ReadFile("sql/sql_files.json")
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(c, sf)
	if err != nil {
		log.Fatalln(err)
	}
	sqlStr := "select count(1) from db_upgrade_log where name=?"
	var count int
	for _, filename := range sf.Files {
		err = DB.Get(&count, sqlStr, filename)
		if count == 0 {
			_, err = sqlx.LoadFile(DB, filepath.Join("sql", filename))
			if err != nil {
				log.Fatalln(err)
			} else {
				_, err = DB.Exec("INSERT INTO db_upgrade_log (name,exectime) VALUES(?,?)", filename, time.Now().Format("2006.01.02-15:04:05"))
				if err != nil {
					log.Fatalln(err)
				}
				log.Infoln("Upgrade sql file: ", filename)
			}
		}
	}
}
