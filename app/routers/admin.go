package routers

import (
	"github.com/gin-gonic/gin"
	"ws/app/auth"
	http "ws/app/http/handlers/admin"
	middleware "ws/app/http/middleware/admin"
	"ws/app/util"
	"ws/app/websocket"
)

func registerAdmin() {
	g := Router.Group("/backend")
	g.POST("/login", http.Login)
	authGroup := g.Group("/")
	authGroup.Use(middleware.Authenticate)
	authGroup.GET("/me", http.Me)
	authGroup.POST("/me/avatar", http.Avatar)
	authGroup.GET("/me/settings" , http.GetChatSetting)
	authGroup.PUT("/me/settings" , http.UpdateChatSetting)
	authGroup.POST("/me/settings/image", http.ChatSettingImage)
	authGroup.DELETE("/ws/chat-user/:id", http.RemoveUser)
	authGroup.POST("/ws/chat-user", http.AcceptUser)
	authGroup.GET("/ws/chat-users", http.ChatUserList)
	authGroup.POST("/ws/read-all", http.ReadAll)
	authGroup.POST("/ws/image", http.Image)
	authGroup.GET("/ws/messages", http.GetHistoryMessage)
	authGroup.GET("/ws/user/:id", http.GetUserInfo)
	authGroup.GET("/ws/sessions/:uid", http.GetHistorySession)

	authGroup.POST("/ws/transfer", http.Transfer)
	authGroup.GET("/ws/transfer/:id/messages", http.TransferMessages)


	authGroup.GET("/settings", http.GetSettings)
	authGroup.PUT("/settings/:name", http.UpdateSetting)

	authGroup.GET("/auto-messages", http.GetAutoMessages)
	authGroup.POST("/auto-messages", http.StoreAutoMessage)
	authGroup.PUT("/auto-messages/:id", http.UpdateAutoMessage)
	authGroup.DELETE("/auto-messages/:id", http.DeleteAutoMessage)
	authGroup.GET("/auto-messages/:id", http.ShowAutoMessage)
	authGroup.POST("/auto-messages/image", http.StoreAutoMessageImage)

	authGroup.GET("/system-auto-rules", http.GetSystemRules)
	authGroup.PUT("/system-auto-rules", http.UpdateSystemRules)

	authGroup.GET("/options/messages", http.GetSelectAutoMessage)
	authGroup.GET("/options/scenes", http.GetSelectScene)
	authGroup.GET("/options/events", http.GetSelectEvent)

	authGroup.POST("/auto-rules", http.StoreAutoRule)
	authGroup.PUT("/auto-rules/:id", http.UpdateAutoRule)
	authGroup.GET("/auto-rules", http.GetAutoRules)
	authGroup.GET("/auto-rules/:id", http.ShowAutoRule)
	authGroup.DELETE("/auto-rules/:id", http.DeleteAutoRule)



	authGroup.GET("/chat-sessions", http.GetChatSession)
	authGroup.GET("/chat-sessions/:id", http.GetChatSessionDetail)

	authGroup.GET("/dashboard/query-info", http.GetUserQueryInfo)
	authGroup.GET("/dashboard/online-info", http.GetOnlineInfo)


	authGroup.GET("/ws", func(c *gin.Context) {
		serviceUser := auth.GetAdmin(c)
		conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			util.RespError(c , err.Error())
			return
		}
		client := websocket.NewAdminConn(serviceUser, conn)
		client.Setup()
		websocket.AdminHub.Login(client)
	})
}
