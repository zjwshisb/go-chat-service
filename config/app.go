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
	Url string
}
type redis struct {
	Addr string
	Auth string
}
type app struct {
	StoragePath string
	LogPath string
}
var (
	Mysql = &mysql{}
	Http  = &http{}
	Redis = &redis{}
	App = &app{}
)
func Setup() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatal(err)
	}
	err = cfg.Section("Mysql").MapTo(Mysql)
	if err != nil {
		log.Fatal(err)
	}
	err = cfg.Section("Http").MapTo(Http)
	if err != nil {
		log.Fatal(err)
	}
	err = cfg.Section("Redis").MapTo(Redis)
	if err != nil {
		log.Fatal(err)
	}
	err = cfg.Section("App").MapTo(App)
	if err != nil {
		log.Fatal(err)
	}
}
