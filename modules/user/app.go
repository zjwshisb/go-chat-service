package user

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	sHttp "ws/modules/service/http"
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
	g := routers.Router.Group("/user")
	{
		g.POST("/login", sHttp.Login)

		g.GET("/ws", func(c *gin.Context) {
			conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				return
			}
			clint := &UClient{
				Conn: conn,
				Send: make(chan []byte, 1000),
				UserId: 1,
			}
			go clint.GetMsg()
			go clint.SendMsg()
			go clint.ReadMsg()
		})
	}

}
