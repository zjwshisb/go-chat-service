package main

import (
	"github.com/sirupsen/logrus"
	"os"
	"time"
	"ws/config"
)
var (
	logFile *os.File
)
func init()  {
	logPath := config.Config["LOG_PATH"]
	fileInfo, err := os.Stat(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(logPath, 0666)
		}
	} else {
		if !fileInfo.IsDir() {
			os.Mkdir(logPath, 0666)
		}
	}
	logName := time.Now().Format("2006-01-02") + ".txt"
	logFile, err := os.OpenFile(logPath + "/" + logName, os.O_CREATE | os.O_WRONLY | os.O_APPEND , 0666)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(logFile)
}
