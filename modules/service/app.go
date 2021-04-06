package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
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
		auth.GET("me", sHttp.Me)

		auth.GET("/ws", func(c *gin.Context) {
			u, _ := c.Get("user")
			serverUser := u.(*models.ServerUser)
			conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				fmt.Println(err)
			}
			client := &hub.Client{
				Conn:        conn,
				UserId:      serverUser.ID,
				Send:        make(chan *models.Action, 1000),
				CloseSignal: make(chan struct{}),
			}
			hub.Hub.Server.Login(client)
		})
	}
}
