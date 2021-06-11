package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	http "ws/internal/http/service"
	middleware "ws/internal/middleware/service"
	"ws/internal/models"
	"ws/internal/websocket"
)

func registerService()  {
	g := Router.Group("/service")
	{
		g.POST("/login", http.Login)

		auth := g.Group("/")

		auth.Use(middleware.Authenticate)
		auth.GET("/me", http.Me)
		auth.POST("/me/avatar", http.Avatar)
		auth.DELETE("/ws/chat-user/:id", http.RemoveUser)
		auth.POST("/ws/chat-user", http.AcceptUser)
		auth.GET("/ws/chat-users", http.ChatUserList)
		auth.POST("/ws/read-all", http.ReadAll)
		auth.POST("/ws/image", http.Image)
		auth.GET("/ws/messages", http.GetHistoryMessage)

		auth.GET("/replies", http.GetShortcutReply)
		auth.POST("/replies", http.StoreShortcutReply)
		auth.PUT("/replies/:id", http.UpdateShortcutReply)
		auth.DELETE("replies/:id", http.DeleteShortcutReply)

		auth.GET("/ws", func(c *gin.Context) {
			ui, _ := c.Get("user")
			serviceUser := ui.(*models.ServiceUser)
			conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				fmt.Println(err)
			}
			client := websocket.NewServiceConn(serviceUser, conn)
			client.Setup()
			websocket.ServiceHub.Login(client)
		})
	}
}
