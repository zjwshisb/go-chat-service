package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"ws/message"
)

func main() {
	var upgrade = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	var connMap = make(map[*Client]bool)
	var channel = make(chan []byte, 1000)
	server := gin.Default()
	server.GET("/", func(context *gin.Context) {
	})
	server.GET("/ws", func(context *gin.Context) {
		var conn, err = upgrade.Upgrade(context.Writer, context.Request, nil)
		fmt.Print(err)
		if err == nil {
			go func(conn *websocket.Conn) {
				client := NewClient(conn)
				connMap[client] = true
				defer func() {
					conn.Close()
					connMap[client] = false
				}()
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
	go func() {
		for {
			msg := <-channel
			for con, exist := range connMap {
				if exist {
					con.Conn.WriteJSON(message.Message{MType: "1", Name: string(msg)})
				}
			}
		}
	}()
	server.Run(":5000")
}
