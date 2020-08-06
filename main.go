package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"ws/uclient"
)

func main() {
	var upgrade = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	var channel = make(chan []byte, 1000)
	server := gin.Default()
	server.GET("/", func(context *gin.Context) {
	})
	server.GET("/ws", func(context *gin.Context) {
		var conn, err = upgrade.Upgrade(context.Writer, context.Request, nil)
		fmt.Print(err)
		if err == nil {
			go func(conn *websocket.Conn) {
				client := uclient.NewClient(conn, 1)
				uclient.WaitAccepts[client] = true
				fmt.Print("链接开始\n")
				for {
					mType, msg, err := conn.ReadMessage()
					if err != nil {
						fmt.Print(err)
						break
					}
					conn.WriteMessage(mType, []byte(string(rune(client.Id))))
					channel <- msg
				}
			}(conn)
		}
	})
	server.GET("/index", func(context *gin.Context) {
		context.File("index.html")
	})
	server.Run(":5000")
}
