package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"ws/app"
	"ws/app/routers"
	"ws/configs"
)
const (
	pidName = "./ws.pid"
)

var pidFile *os.File

func init() {
	pidFile, _ = os.OpenFile(pidName, os.O_WRONLY | os.O_CREATE, 0755 )
	if err := syscall.Flock(int(pidFile.Fd()), syscall.LOCK_EX | syscall.LOCK_NB); err != nil {
		log.Fatalln("server is running ...")
	}
	pid := os.Getpid()
	_, err := pidFile.Write([]byte(strconv.Itoa(pid)))
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	app.Setup()
	srv := &http.Server{
		Addr:    configs.Http.Host +":" + configs.Http.Port,
		Handler: routers.Router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-quit
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer func() {
		cancel()
		if err := syscall.Flock(int(pidFile.Fd()), syscall.LOCK_UN); err != nil {
			log.Fatalln(err)
		}
	}()
	if err:= srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
