package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"ws/app/util"
)

var Log *logrus.Logger

func Setup()  {


	env := util.GetEnv()

	var level logrus.Level

	switch strings.ToLower(env) {
	case "production":
		level = logrus.ErrorLevel
	case "test":
		level = logrus.DebugLevel
	case "local":
		level = logrus.DebugLevel
	}

	Log = logrus.New()
	Log.SetLevel(level)
	Log.SetReportCaller(true)
	Log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat:"2006-01-02 15:04:05",
	})
	Log.SetFormatter(&logrus.JSONFormatter{})
	LogFile := viper.GetString("App.LogFile")
	path, _ := filepath.Split(LogFile)
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			panic(fmt.Errorf("log error: %s is not exist", path ))
		}
	}
	write, err := os.OpenFile(LogFile, os.O_APPEND | os.O_WRONLY | os.O_CREATE, 0755 )
	if err != nil {
		panic(fmt.Errorf("log err: %v", err ))
	}
	Log.SetOutput(write)
}