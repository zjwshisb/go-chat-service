package routers

import (
	"github.com/gin-gonic/gin"
	"ws/internal/http/user"
	user2 "ws/internal/middleware/user"
	"ws/internal/models"
	hub "ws/internal/websocket"
)

func registerUser()  {
	u := Router.Group("/user")
	{
		u.POST("/login", user.Login)
		auth := u.Group("/")
		auth.Use(user2.Authenticate)
		auth.GET("/ws/messages", user.GetHistoryMessage)
		auth.POST("/ws/image", user.Image)
		auth.GET("/ws", func(c *gin.Context) {
			conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				return
			}
			ui, _ := c.Get("user")
			userModel := ui.(*models.User)
			client := hub.NewUserConn(userModel, conn)
			client.Setup()
			hub.UserHub.Login(client)
		})
	}
}
