package configs

import (
	"flag"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"log"
	"os"
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
	Db int
}
type app struct {
	LogPath string
	LogLevel logrus.Level
	Url string
	Env string
	SystemChatName string
	PidFile string
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

	var file string

	flag.StringVar(&file,"c", "config.ini", "config file")
	flag.Parse()


	_, err := os.Stat(file)
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := ini.Load(file)

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
