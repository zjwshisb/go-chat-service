package configs

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"log"
)

type wechat struct {
	MiniProgramAppId string
	MiniProgramAppSecret string
	SubscribeTemplateIdOne string
}
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
type app struct {
	LogPath string
	LogLevel logrus.Level
	Url string
}
type file struct {
	Storage string
	LocalPath string
	LocalPrefix string
	QiniuAk string
	QiniuSK string
	QiniuUrl string
	QiniuBucket string
}
var (
	Mysql = &mysql{}
	Http  = &http{}
	Redis = &redis{}
	App = &app{}
	File = &file{}
	Wechat = &wechat{}
)
func init() {

	cfg, err := ini.Load("../../../config.ini")

	err = cfg.Section("Mysql").MapTo(Mysql)

	err = cfg.Section("Http").MapTo(Http)

	err = cfg.Section("Redis").MapTo(Redis)

	err = cfg.Section("App").MapTo(App)

	err = cfg.Section("File").MapTo(File)

	err = cfg.Section("Wechat").MapTo(Wechat)
	if err != nil {
		log.Fatal(err)
	}
}
