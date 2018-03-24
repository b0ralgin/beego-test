package utilities

import (
	"errors"

	"github.com/astaxie/beego"
	"github.com/globalsign/mgo"
)

var (
	singleton *mgo.Session
)

func Startup() (err error) {
	url := beego.AppConfig.String("db_url")
	if url == "" {
		return errors.New("can't find Mongo URL")
	}
	if singleton == nil {
		singleton, err = mgo.Dial(url)
		return err
	}
	return singleton.Ping()
}

func Close() {
	singleton.Close()
}

func CopySession() *mgo.Session {
	return singleton.Clone()
}
