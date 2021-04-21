package user

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
	"ws/action"
	"ws/hub"
	"ws/models"
	sHttp "ws/modules/user/http"
	"ws/modules/user/middleware"
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
		auth := g.Group("/")
		auth.Use(middleware.Authenticate)
		auth.GET("/ws/messages",sHttp.GetHistoryMessage)
		auth.GET("/ws", func(c *gin.Context) {
			conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				return
			}
			u, _ := c.Get("user")
			user := u.(*models.User)
			client := &hub.UClient{
				Conn: conn,
				Send: make(chan *action.Action, 1000),
				User: user,
				CloseSignal: make(chan struct{}),
				CreatedAt: time.Now().Unix(),
			}
			client.Setup()
			hub.Hub.User.Login(client)
		})
	}

}
