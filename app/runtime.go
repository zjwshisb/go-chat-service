package app

import (
	"log"
	"os"
	"strconv"
	"syscall"
)

var (
	pidName = "./ws.pid"
	pidFile *os.File
)

func init()  {
	pidFile, _ = os.OpenFile(pidName, os.O_WRONLY|os.O_CREATE, 0755)
	if err := syscall.Flock(int(pidFile.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		log.Fatalln("server is running ...")
	}
	pidFile.Truncate(0)
	pidFile.Seek(0,0)
	pid := os.Getpid()
	_, err := pidFile.Write([]byte(strconv.Itoa(pid)))
	if err != nil {
		log.Fatalln(err)
	}
}

func Clear() {
	if err := syscall.Flock(int(pidFile.Fd()), syscall.LOCK_UN); err != nil {
		log.Fatalln(err)
	}
}

