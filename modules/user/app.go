package user

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
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
		auth.POST("/ws/image", sHttp.Image)
		auth.GET("/ws", func(c *gin.Context) {
			conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				return
			}
			u, _ := c.Get("user")
			user := u.(*models.User)
			client := hub.NewUserConn(user, conn)
			client.Setup()
			hub.UserHub.Login(client)
		})
	}

}
