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
	"syscall"
	"time"
	"ws/app/cron"
	_ "ws/app/http/requests"
	_ "ws/app/log"
	"ws/app/routers"
	"ws/app/websocket"
)

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


