package app

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
	"ws/app/routers"
	"ws/app/websocket"
	"ws/configs"
)

func getLogPid()  int {
	pidName := configs.App.PidFile
	pidFile, err := os.Open(pidName)
	if err != nil {
		log.Fatalln(err)
	}
	bytes , err := ioutil.ReadAll(pidFile)
	if err != nil {
		pidFile.Close()
		log.Fatalln(err)
	}
	str := string(bytes)
	pid, err := strconv.Atoi(str)
	if err != nil {
		pidFile.Close()
		log.Fatalln(err)
	}
	pidFile.Close()
	return pid
}

func logPid(pid int)  {
	pidName := configs.App.PidFile
	pidFile, _ := os.OpenFile(pidName, os.O_WRONLY|os.O_CREATE, 0755)
	pidFile.Truncate(0)
	pidFile.Seek(0,0)
	_, err := pidFile.Write([]byte(strconv.Itoa(pid)))
	if err != nil {
		pidFile.Close()
		log.Fatalln(err)
	}
	pidFile.Close()
}

func Start()  {
	quit := make(chan os.Signal, 1)
	websocket.Setup()
	routers.Setup()
	srv := &http.Server{
		Addr:    configs.Http.Host +":" + configs.Http.Port,
		Handler: routers.Router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			quit<-syscall.SIGKILL
			log.Fatalln(err)
		}
	}()
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	logPid(os.Getpid())
	fmt.Println("server start")
	<-quit
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer func() {
		cancel()
	}()
	if err:= srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exited")
}

func Stop()  {
	pid := getLogPid()
	cmd := exec.Command("ps", "aux")
	output, _ := cmd.Output()
	index := strings.Index(string(output), strconv.Itoa(pid))
	if index > 0 {
		closeCmd := exec.Command("kill", "-2" , strconv.Itoa(pid))
		result, _ := closeCmd.Output()
		fmt.Println(string(result))
	} else {
		log.Fatalln("server is not runing")
	}
}

