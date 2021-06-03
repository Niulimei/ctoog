package restapi

import (
	"ctgb/utils"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func CheckPermission(r *http.Request) bool {
	username := r.Header.Get("username")
	userInfo := getUserInfo(username)
	if userInfo.RoleID != int64(AdminRole) {
		utils.RecordLog(utils.Info, utils.Auth, "", "user "+username+" Forbidden", http.StatusForbidden)
		return false
	}
	return true
}

var exceptionURL = map[string]string{
	"/api/login":      "POST",
	"/users/register": "POST",
}

func IsExceptionURL(method, uri string) bool {
	for uriChild, m := range exceptionURL {
		if strings.HasSuffix(uri, uriChild) && m == method {
			return true
		}
	}
	return false
}

func DumpLogFile(logFile string) {
	for {
		N := time.Now()
		y, m, d := N.Date()
		delay := time.Date(y, m, d+1, 0, 0, 0, 0, time.Local).Unix() - N.Unix()
		//delay := time.Date(y, m, d, N.Hour(), N.Minute()+1, 0, 0, time.Local).Unix() - N.Unix()
		time.Sleep(time.Second * time.Duration(delay))
		log.Debug("start bak")
		fs, err := ioutil.ReadDir(filepath.Dir(logFile))
		if err != nil {
			log.Debug(err)
		}
		for _, f := range fs {
			if time.Now().Unix()-f.ModTime().Unix() > 3*24*3600 {
				//if time.Now().Minute()-f.ModTime().Minute() > 1 && filepath.Base(logFile) != f.Name() {
				log.Debug("start clean")
				err = os.RemoveAll(filepath.Join(filepath.Dir(logFile), f.Name()))
				if err != nil {
					log.Debug(err)
				}
			}
		}
		s, err := os.Stat(logFile)
		if err != nil {
			log.Debug(err)
			continue
		}
		co, err := ioutil.ReadFile(logFile)
		if err != nil {
			log.Debug(err)
			continue
		}
		err = ioutil.WriteFile(logFile, []byte(""), s.Mode())
		if err != nil {
			log.Debug(err)
			continue
		}
		err = ioutil.WriteFile(logFile+"_"+time.Now().Format("2006.01.02_15:04:05"), co, s.Mode())
		if err != nil {
			log.Debug(err)
			continue
		}
	}
}
