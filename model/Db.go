package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"ws/config"
)
var (
	Db *gorm.DB
)

func init()  {
	conf := config.DBConf
	host := conf.Host + ":" + conf.Port
	extra := "?charset=utf8&parseTime=True&loc=Local"
	user := conf.Username + ":" + conf.Password
	var err error
	Db, err = gorm.Open(conf.Connection, user + "@/" + host + extra)
	if err != nil {
		log.Fatal(err)
	}
}
