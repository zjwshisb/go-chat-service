package configs

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"log"
	"os"
	"ws/command"
)

type wechat struct {
	MiniProgramAppId string
	MiniProgramAppSecret string
	SubscribeTemplateIdOne string // 订阅消息模板id
	ChatPath string // 小程序客户页面路径
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
	SystemChatName string // 系统消息客服名称
	PidFile string  // pid文件
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
	_, err := os.Stat(command.ConfigFile)
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := ini.Load(command.ConfigFile)

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
