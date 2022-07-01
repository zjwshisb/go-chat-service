package sys

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/duke-git/lancet/v2/netutil"
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
func GetPidFile() string {
	ips := netutil.GetIps()
	s := []byte(ips[0])
	h := md5.New()
	h.Write(s)
	return hex.EncodeToString(h.Sum(nil)) + ".pid"
}

func GetPid() int {
	dir := config.GetStoragePath()
	pidFile := dir + "/" + GetPidFile()
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
	pidFile := dir + "/" + GetPidFile()
	file, err := os.OpenFile(pidFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("pid file err: %v", err)
	}
	defer file.Close()
	file.Write([]byte(strconv.Itoa(pid)))
}
