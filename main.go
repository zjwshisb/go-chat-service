package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
	"ws/uclient"
	"ws/database"
)

var (
	action string
)

func main() {
	flag.StringVar(&action, "a", "run", "set run type")
	flag.Parse()
	if action == "migrate" {
		database.Migrate()
	} else {
		var upgrade = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		server := gin.Default()
		server.GET("/", func(context *gin.Context) {
		})
		server.GET("/ws", func(context *gin.Context) {
			var conn, err = upgrade.Upgrade(context.Writer, context.Request, nil)
			fmt.Print(err)
			if err == nil {
				go func(conn *websocket.Conn) {
					client := uclient.NewClient(conn, 1)
					go func() {
						ticker := time.NewTicker(time.Second * 1)
						for  {
							select {
								case <- ticker.C:
									client.Push([]byte("1"))
							}
						}
					}()
					fmt.Println(client)
				}(conn)
			}
		})
		server.GET("/index", func(context *gin.Context) {
			context.File("index.html")
		})
		server.Run(":5000")
	}
}
