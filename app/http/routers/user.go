package routers

import (
	http "ws/app/http/controllers/user"
	middleware "ws/app/http/middleware/user"
	"ws/app/http/websocket"
	"ws/app/models"

	"github.com/gin-gonic/gin"
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
		auth.POST("/ws/read", http.ReadAll)
		auth.GET("/ws", func(c *gin.Context) {
			conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				return
			}
			ui, _ := c.Get("frontend")
			userModel := ui.(*models.User)
			client := websocket.NewConn(userModel, conn, websocket.UserManager)
			websocket.UserManager.Register(client)
		})
	}
}
