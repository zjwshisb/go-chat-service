package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"ws/action"
	"ws/hub"
	"ws/models"
	sHttp "ws/modules/service/http"
	"ws/modules/service/middleware"
	"ws/routers"
)

var (
	upgrade = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func Setup() {
	g := routers.Router.Group("/service")
	{
		g.POST("/login", sHttp.Login)

		auth := g.Group("/")

		auth.Use(middleware.Authenticate)
		auth.GET("/me", sHttp.Me)
		auth.POST("/me/avatar", sHttp.Avatar)
		auth.DELETE("/ws/chat-user/:id", sHttp.RemoveUser)
		auth.POST("/ws/chat-user", sHttp.AcceptUser)
		auth.GET("/ws/chat-users", sHttp.ChatUserList)
		auth.POST("/ws/read-all", sHttp.ReadAll)
		auth.POST("/ws/image", sHttp.Image)
		auth.GET("/ws/messages", sHttp.GetHistoryMessage)

		auth.GET("/ws", func(c *gin.Context) {
			u, _ := c.Get("user")
			serverUser := u.(*models.ServerUser)
			conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				fmt.Println(err)
			}
			client := &hub.Client{
				Conn:        conn,
				User:      serverUser,
				Send:        make(chan *action.Action, 1000),
				CloseSignal: make(chan struct{}),
			}
			hub.Hub.Server.Login(client)
		})

	}
}
