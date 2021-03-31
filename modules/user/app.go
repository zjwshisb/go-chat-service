package user

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"ws/routers"
	"ws/util"
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

		g.GET("/ws", func(c *gin.Context) {
			conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				return
			}
			clint := &UClient{
				Conn: conn,
				Send: make(chan *util.Action, 1000),
				UserId: 1,
			}
			go clint.GetMsg()
			go clint.SendMsg()
			go clint.ReadMsg()
		})
	}

}
