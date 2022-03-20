package routers

import (
	"github.com/gin-gonic/gin"
	http "ws/app/http/controllers/user"
	middleware "ws/app/http/middleware/user"
	"ws/app/http/websocket"
	"ws/app/models"
)

func registerFrontend() {
	u := Router.Group("/user")
	{
		u.POST("/login", http.Login)
		auth := u.Group("/")
		auth.Use(middleware.Authenticate)
		auth.GET("/template-id", http.GetTemplateId)
		auth.POST("/subscribe", http.Subscribe)
		auth.GET("/ws/messages", http.GetHistoryMessage)
		auth.POST("/ws/image", http.Image)
		auth.POST("/ws/req-id", http.GetReqId)
		auth.GET("/ws", func(c *gin.Context) {
			conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				return
			}
			ui, _ := c.Get("frontend")
			userModel := ui.(*models.User)
			client := websocket.NewUserConn(userModel, conn)
			websocket.UserManager.Register(client)
		})
	}
}
