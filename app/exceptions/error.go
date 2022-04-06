package exceptions

import (
	"runtime"
	"ws/app/log"
)

func Handler(err error) {
	_, file, line, _ := runtime.Caller(1)
	log.Log.Error(file, line, err)
}
