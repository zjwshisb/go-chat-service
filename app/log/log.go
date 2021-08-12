package log

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"time"
	"ws/configs"
)

var Log *logrus.Logger

func init()  {
	Log = logrus.New()
	Log.SetLevel(configs.App.LogLevel)
	Log.SetReportCaller(true)
	logPath := configs.App.LogPath
	_, err := os.Stat(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatal(logPath + " is not exist")
		}
	}
	filename := time.Now().Format("2006-01-02") + ".log"
	write, err := os.OpenFile(logPath + "/" + filename, os.O_WRONLY | os.O_CREATE, 0755 )
	if err != nil {
		log.Fatal(err)
	}
	Log.SetOutput(write)
}