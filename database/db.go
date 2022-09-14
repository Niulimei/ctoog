package database

import (
	"ctgb/models"
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

type TaskModel struct {
	// id
	Id models.JsonNullInt64 `db:"id"`

	// cc password
	CcPassword models.JsonNullString `db:"cc_password"`

	// cc user
	CcUser models.JsonNullString `db:"cc_user"`

	// component
	Component models.JsonNullString

	// dir
	Dir models.JsonNullString

	// keep
	Keep models.JsonNullString

	// git password
	GitPassword models.JsonNullString `db:"git_password"`

	// git URL
	GitURL models.JsonNullString `db:"git_url"`

	// git user
	GitUser models.JsonNullString `db:"git_user"`

	// status
	Status models.JsonNullString

	//last completed date time
	LastCompletedDateTime models.JsonNullString `db:"last_completed_date_time"`

	//worker id
	WorkerId models.JsonNullInt64 `db:"worker_id"`

	// pvob
	Pvob models.JsonNullString

	// include empty dir
	IncludeEmpty models.JsonNullBool `db:"include_empty"`

	// git email
	GitEmail models.JsonNullString `db:"git_email"`

	// svn url
	SvnURL models.JsonNullString `db:"svn_url"`

	// model type
	ModelType models.JsonNullString `db:"model_type"`

	// gitignore
	Gitignore models.JsonNullString `json:"gitignore,omitempty"`

	// branches_info
	BranchesInfo models.JsonNullString `json:"branches_info,omitempty" db:"branches_info"`

	// gitlab group
	GitlabGroup models.JsonNullString `json:"gitlab_group" db:"gitlab_group"`

	// gitlab project
	GitlabProject models.JsonNullString `json:"gitlab_project" db:"gitlab_project"`

	// gitlab token
	GitlabToken models.JsonNullString `json:"gitlab_token" db:"gitlab_token"`

	// gitee group
	GiteeGroup models.JsonNullString `json:"gitee_group" db:"gitee_group"`

	// gitee project
	GiteeProject models.JsonNullString `json:"gitee_project" db:"gitee_project"`

	// gitee token
	GiteeToken models.JsonNullString `json:"gitee_token" db:"gitee_token"`

	// source url
	SourceURL models.JsonNullString `json:"source_url" db:"source_url"`

	// target url
	TargetURL models.JsonNullString `json:"target_url" db:"target_url"`
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

type PlanModel struct {
	ID               int64
	Status           string
	OriginType       string `db:"origin_type"`
	Pvob             string
	Component        string
	Dir              string
	OriginURL        string `db:"origin_url"`
	TranslateType    string `db:"translate_type"`
	TargetURL        string `db:"target_url"`
	Subsystem        string
	ConfigLib        string `db:"config_lib"`
	Group            string `db:"business_group"`
	Team             string
	Supporter        string
	SupporterTel     string `db:"supporter_tel"`
	Creator          string
	Tip              string
	ProjectType      string `db:"project_type"`
	Purpose          string
	PlanStartTime    string `db:"plan_start_time"`
	PlanSwitchTime   string `db:"plan_switch_time"`
	ActualStartTime  string `db:"actual_start_time"`
	ActualSwitchTime string `db:"actual_switch_time"`
	Effect           string
	TaskID           int64 `db:"task_id"`
	Extra1           string
	Extra2           string
	Extra3           string
}

func InitializeDataBase() {
	initDB()
	upgradeDB()
	go func() {
		for {
			startTime := time.Now().Format("2006.01.02-15:04:05")
			cmd := exec.Command("cp", "store/translator.db", "backup/translator-"+startTime+".back")
			cmd.Run()
			time.Sleep(time.Hour * 24)
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
	_, err = os.Stat("store/translator.db") //os.Stat获取文件信息
	if err != nil {
		if os.IsNotExist(err) {
			isInitAlready = false
		}
	}
	DB, err = sqlx.Connect("sqlite3", "file:store/translator.db?cache=private&mode=rwc")
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
