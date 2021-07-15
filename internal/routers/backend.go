package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"ws/internal/auth"
	http "ws/internal/http/handlers/backend"
	middleware "ws/internal/http/middleware/backend"
	"ws/internal/websocket"
)

func registerBackend() {
	g := Router.Group("/backend")
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

	//authGroup.GET("/user/:id", http.GetUserInfo)

	authGroup.GET("/settings", http.GetSettings)
	authGroup.PUT("/settings/:name", http.UpdateSetting)

	authGroup.GET("/auto-messages", http.GetAutoMessages)
	authGroup.POST("/auto-message", http.StoreAutoMessage)
	authGroup.PUT("/auto-message/:id", http.UpdateAutoMessage)
	authGroup.DELETE("/auto-message/:id", http.DeleteAutoMessage)
	authGroup.GET("/auto-message/:id", http.ShowAutoMessage)
	authGroup.POST("/auto-message/image", http.StoreAutoMessageImage)

	authGroup.GET("/system-auto-rules", http.GetSystemRules)
	authGroup.PUT("/system-auto-rules", http.UpdateSystemRules)

	authGroup.GET("/auto-rules/options/messages", http.GetSelectAutoMessage)
	authGroup.POST("/auto-rule", http.StoreAutoRule)
	authGroup.GET("/auto-rule/:id", http.ShowAutoRule)
	authGroup.GET("/auto-rules", http.GetAutoRules)
	authGroup.PUT("/auto-rule/:id", http.UpdateAutoRule)

	authGroup.GET("/dashboard/query-info", http.GetUserQueryInfo)

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
