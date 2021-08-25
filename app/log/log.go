package log

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path/filepath"
	"ws/configs"
)

var Log *logrus.Logger

func init()  {
	Log = logrus.New()
	Log.SetLevel(configs.App.LogLevel)
	Log.SetReportCaller(true)
	LogFile := configs.App.LogFile
	path, filename := filepath.Split(LogFile)
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatal(path + " is not exist")
		}
	}
	write, err := os.OpenFile(path + "/" + filename, os.O_WRONLY | os.O_CREATE, 0755 )
	if err != nil {
		log.Fatal(err)
	}
	Log.SetOutput(write)
}