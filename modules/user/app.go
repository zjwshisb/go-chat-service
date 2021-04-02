package user

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"ws/hub"
	"ws/models"
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

		g.GET("/ws", func(c *gin.Context) {
			conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				return
			}
			client := &hub.UClient{
				Conn: conn,
				Send: make(chan *models.Action, 1000),
				UserId: 1,
				CloseSignal: make(chan struct{}),
			}
			hub.Hub.User.Login(client)
		})
	}

}
