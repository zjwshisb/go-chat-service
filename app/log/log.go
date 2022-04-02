package log

import (
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"ws/app/util"
	"ws/config"
)

type log struct {
	*logrus.Logger
}

var Log *log

var prefix = "log"

func Setup() {
	env := config.GetEnv()
	var level logrus.Level
	var output io.Writer
	switch strings.ToLower(env) {
	case "production":
		fallthrough
	case "test":
		level = logrus.ErrorLevel
		storagePath := config.GetStoragePath()
		logPath := storagePath + "/" + prefix
		if !util.DirExist(logPath) {
			err := util.MkDir(logPath)
			if err != nil {
				panic(fmt.Sprintf("make log dir[%s] err: %v", logPath, err))
			}
		}
		LogFile := logPath + "/" + "app.log"
		var err error
		output, err = os.OpenFile(LogFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			panic(fmt.Errorf("log file err: %v", err))
		}
	case "local":
		fallthrough
	default:
		output = os.Stdout
		level = logrus.DebugLevel
	}
	Log = &log{
		Logger: logrus.New(),
	}
	Log.Logger.SetLevel(level)
	Log.Logger.SetReportCaller(true)
	Log.Logger.SetFormatter(&nested.Formatter{
		TimestampFormat: "2006/01/02 03:04:05",
		NoFieldsSpace:   true,
		HideKeys:        true,
		CustomCallerFormatter: func(r *runtime.Frame) string {
			return " " + r.File + ":" + strconv.Itoa(r.Line)
		},
	})
	Log.Logger.SetOutput(output)
}
