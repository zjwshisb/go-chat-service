package configs

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"log"
	"net"
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
	LogFile string
	LogLevel logrus.Level
	Url string
	Env string
	SystemChatName string // 系统消息客服名称
	PidFile string  // pid文件
	PublicPath string
	Name string
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
	if cfg != nil {
		_ = cfg.Section("Mysql").MapTo(Mysql)
		_ = cfg.Section("Http").MapTo(Http)
		_ = cfg.Section("Redis").MapTo(Redis)
		_ = cfg.Section("App").MapTo(App)
		ips, err := getLocalIP()
		if err != nil {
			log.Fatalln(err)
		}
		App.Name = ips[0] + ":" + Http.Port
		err = cfg.Section("File").MapTo(File)

		err = cfg.Section("Wechat").MapTo(Wechat)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func getLocalIP() (ips []string, err error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, i := range ifaces {
		addrs, errRet := i.Addrs()
		if errRet != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
				if ip.IsGlobalUnicast() {
					ips = append(ips, ip.String())
				}
			}
		}
	}
	return
}
