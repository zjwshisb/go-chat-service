package app

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"ws/app/cron"
	_ "ws/app/http/requests"
	_ "ws/app/log"
	"ws/app/routers"
	"ws/app/util"
	"ws/app/websocket"
)

func GetPid() int {
	dir := util.GetWorkDir()
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
	dir := util.GetWorkDir()
	pidFile := dir + "/pid.log"
	file, err := os.OpenFile(pidFile, os.O_CREATE | os.O_TRUNC | os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
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


