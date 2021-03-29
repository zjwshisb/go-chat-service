package config

import (
	"gopkg.in/ini.v1"
	"log"
)

type mysql struct {
	Username string
	Password string
	Host string
	Name string
	Port string
}
type http struct {
	Port string
	Host string
}
type redis struct {
	Addr string
	Auth string
}

var (
	Mysql *mysql
	Http *http
	Redis *redis
)
func Setup() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatal(err)
	}
	Mysql = &mysql{}
	err = cfg.Section("Mysql").MapTo(Mysql)
	Http = &http{}
	err = cfg.Section("Http").MapTo(Http)
	Redis = &redis{}
	err = cfg.Section("Redis").MapTo(Redis)
	if err != nil {
		log.Fatal(err)
	}
}
