package sys

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"ws/config"
)

func IsRunning() bool {
	pid := GetPid()
	if pid == 0 {
		return false
	} else {
		cmd := exec.Command("ps")
		out, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}
		return strings.Contains(string(out), strconv.Itoa(pid))
	}
}

func GetPid() int {
	dir := config.GetStoragePath()
	pidFile := dir + "/pid.log"
	b, err := os.ReadFile(pidFile)
	if err != nil {
		return 0
	}
	s := string(b)
	pid, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return pid
}

func LogPid() {
	pid := os.Getpid()
	dir := config.GetStoragePath()
	pidFile := dir + "/pid.log"
	file, err := os.OpenFile(pidFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("pid file err: %v", err)
	}
	defer file.Close()
	file.Write([]byte(strconv.Itoa(pid)))
}
