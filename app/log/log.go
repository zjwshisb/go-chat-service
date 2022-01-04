package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"ws/app/util"
)

var Log *logrus.Logger

var prefix = "log"

func Setup()  {

	env := util.GetEnv()

	var level logrus.Level

	switch strings.ToLower(env) {
	case "production":
		level = logrus.ErrorLevel
	case "test":
		fallthrough
	case "local":
		fallthrough
	default:
		level = logrus.DebugLevel
	}

	Log = logrus.New()
	Log.SetLevel(level)
	Log.SetReportCaller(true)
	Log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat:"2006-01-02 15:04:05",
	})
	Log.SetFormatter(&logrus.JSONFormatter{})
	storagePath := util.GetStoragePath()
	logPath := storagePath + "/" + prefix
	if !util.DirExist(logPath) {
		err := util.MkDir(logPath)
		if err != nil {
			panic(fmt.Sprintf("make log dir[%s] err: %v", logPath, err))
		}
	}
	LogFile := logPath + "/" + "app.log"
	write, err := os.OpenFile(LogFile, os.O_APPEND | os.O_WRONLY | os.O_CREATE, 0755 )
	if err != nil {
		panic(fmt.Errorf("log file err: %v", err ))
	}
	Log.SetOutput(write)

}