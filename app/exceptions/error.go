package exceptions

import (
	"runtime"
	"strconv"
	"ws/app/log"
)

func Handler(err error) {
	_, file, line, _ := runtime.Caller(1)
	log.Log.SetReportCaller(false)
	log.Log.Error(err, " ", file, ":", strconv.Itoa(line))
	log.Log.SetReportCaller(true)
}
