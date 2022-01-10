package app

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
	"ws/app/cron"
	_ "ws/app/http/requests"
	"ws/app/routers"
	"ws/app/util"
	"ws/app/websocket"
)

func IsRunning() bool  {
	pid := GetPid()
	if pid == 0 {
		return false
	} else {
		cmd := 	exec.Command("ps")
		out, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}
		return strings.Contains(string(out), strconv.Itoa(pid))
	}
}

func GetPid() int {
	dir := util.GetStoragePath()
	pidFile := dir + "/pid.log"
	b,err := os.ReadFile(pidFile)
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

func logPid() {
	pid := os.Getpid()
	dir := util.GetStoragePath()
	pidFile := dir + "/pid.log"
	file, err := os.OpenFile(pidFile, os.O_CREATE | os.O_TRUNC | os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("pid file err: %v" , err)
	}
	defer file.Close()
	file.Write([]byte(strconv.Itoa(pid)))
}

func Start()  {
	routers.Setup()
	srv := &http.Server{
		Addr:    viper.GetString("Http.Host") +":" +  viper.GetString("Http.Port"),
		Handler: routers.Router,
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		err := srv.ListenAndServe()
		if err != nil  {
			if err != http.ErrServerClosed {
				quit <- syscall.SIGINT
				log.Fatalln(err)
			}
		}
	}()
	defer func() {
		websocket.AdminManager.Destroy()
		websocket.UserManager.Destroy()
		cron.Stop()
	}()
	logPid()
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer func() {
		cancel()
	}()
	if err:= srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown error :", err)
	}
	fmt.Println("exit forced")
}


