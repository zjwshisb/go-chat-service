package log

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"time"
	"ws/configs"
)
var Log *logrus.Logger

func Setup()  {
	Log = logrus.New()
	Log.SetLevel(logrus.WarnLevel)
	Log.SetReportCaller(true)
	logPath := configs.App.LogPath
	filename := time.Now().Format("2006-01-02") + ".log"
	write, err := os.OpenFile(logPath + "/" + filename, os.O_WRONLY | os.O_CREATE, 0755 )
	if err != nil {
		log.Fatal(err)
	}
	Log.Warning("test")
	Log.SetOutput(write)
}