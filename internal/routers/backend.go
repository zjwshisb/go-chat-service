package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"ws/internal/auth"
	http "ws/internal/http/backend"
	middleware "ws/internal/middleware/backend"
	"ws/internal/websocket"
)

func registerBackend()  {
	g := Router.Group("/backend")
	{
		g.POST("/login", http.Login)

		authGroup := g.Group("/")

		authGroup.Use(middleware.Authenticate)
		authGroup.GET("/me", http.Me)
		authGroup.POST("/me/avatar", http.Avatar)
		authGroup.DELETE("/ws/chat-user/:id", http.RemoveUser)
		authGroup.POST("/ws/chat-user", http.AcceptUser)
		authGroup.GET("/ws/chat-users", http.ChatUserList)
		authGroup.POST("/ws/read-all", http.ReadAll)
		authGroup.POST("/ws/image", http.Image)
		authGroup.GET("/ws/messages", http.GetHistoryMessage)

		authGroup.GET("/replies", http.GetShortcutReply)
		authGroup.POST("/replies", http.StoreShortcutReply)
		authGroup.PUT("/replies/:id", http.UpdateShortcutReply)
		authGroup.DELETE("replies/:id", http.DeleteShortcutReply)

		authGroup.GET("/ws", func(c *gin.Context) {
			serviceUser := auth.GetBackendUser(c)
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
