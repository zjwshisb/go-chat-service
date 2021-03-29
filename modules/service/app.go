package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"ws/models"
	"ws/modules"
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
	H *hub
)

func Setup() {
	H = &hub{
		Clients: make(map[int64]*Client),
		Logout: make(chan *Client, 1000),
		Login: make(chan *Client, 1000),
	}
	go H.run()

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
			client := &Client{
				Conn: conn,
				UserId: serverUser.ID,
				isClose: false,
				Send: make(chan *modules.Action, 1000),
				closeSignal: make(chan struct{}),
			}
			client.start()
		})
	}

}
