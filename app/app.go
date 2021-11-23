package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"ws/app/cron"
	_ "ws/app/databases"
	_ "ws/app/http/requests"
	_ "ws/app/log"
	"ws/app/routers"
	"ws/app/websocket"
	"ws/configs"
)

func Start()  {
	routers.Setup()
	srv := &http.Server{
		Addr:    configs.Http.Host +":" + configs.Http.Port,
		Handler: routers.Router,
	}
	quit := make(chan os.Signal, 1)
	go func() {
		err := srv.ListenAndServe()
		if err != nil  {
			if err != http.ErrServerClosed {
				quit<-syscall.SIGINT
				log.Fatalln(err)
			}
		}
	}()
	go func() {
		cron.Run()
	}()
	defer func() {
		websocket.AdminManager.Destroy()
		websocket.UserManager.Destroy()
	}()
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-quit
	log.Println("Shutdown Server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer func() {
		cancel()
	}()
	if err:= srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exited")
}


