package http

import (
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"syscall"
	"ws/app/http/routers"
	"ws/app/http/websocket"
)

func Serve(c chan os.Signal) *http.Server {
	routers.Setup()
	websocket.SetupAdmin()
	websocket.SetupUser()
	srv := &http.Server{
		Addr:    viper.GetString("Http.Host") + ":" + viper.GetString("Http.Port"),
		Handler: routers.Router,
	}
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			if err != http.ErrServerClosed {
				c <- syscall.SIGINT
				fmt.Println(err)
			}
		}
	}()
	return srv
}
