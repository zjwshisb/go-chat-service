package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
	"ws/database"
	"ws/uclient"
)

type app struct {
}

var (
	App *app  = &app{}
	id  int64 = 1
)

func (a *app) migrate() {
	database.Migrate()
}
func (a *app) destroy() {
	logFile.Close()
}
func (a *app) run() {
	defer func() {
		App.destroy()
	}()
	flag.StringVar(&action, "a", "run", "set run type")
	flag.Parse()
	var upgrade = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	server := gin.Default()
	server.GET("/", func(context *gin.Context) {
	})
	server.GET("/ws", func(context *gin.Context) {
		defer func() {
			logrus.Info("connect success")
		}()
		var conn, err = upgrade.Upgrade(context.Writer, context.Request, nil)
		if err != nil {
			logrus.Error(err)
		}
		client := uclient.NewClient(conn, id)
		id++
		go func() {
			ticker := time.NewTicker(time.Second * 1)
			for {
				select {
				case <-ticker.C:
					client.Push([]byte("1"))
				}
			}
		}()
	})
	server.GET("/index", func(context *gin.Context) {
		context.File("index.html")
	})
	server.Run(":5000")
}
func (a *app) seed() {

}
