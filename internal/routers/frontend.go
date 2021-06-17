package routers

import (
	"github.com/gin-gonic/gin"
	"ws/internal/http/frontend"
	middleware "ws/internal/middleware/user"
	"ws/internal/models"
	hub "ws/internal/websocket"
)

func registerFrontend()  {
	u := Router.Group("/frontend")
	{
		u.POST("/login", frontend.Login)
		auth := u.Group("/")
		auth.Use(middleware.Authenticate)
		auth.GET("/template-id", frontend.GetTemplateId)
		auth.GET("/ws/messages", frontend.GetHistoryMessage)
		auth.POST("/ws/image", frontend.Image)
		auth.GET("/ws", func(c *gin.Context) {
			conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				return
			}
			ui, _ := c.Get("frontend")
			userModel := ui.(*models.User)
			client := hub.NewUserConn(userModel, conn)
			client.Setup()
			hub.UserHub.Login(client)
		})
	}
}
