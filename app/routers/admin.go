package routers

import (
	"github.com/gin-gonic/gin"
	"ws/app/auth"
	http "ws/app/http/handlers/admin"
	middleware "ws/app/http/middleware/admin"
	"ws/app/log"
	"ws/app/util"
	"ws/app/websocket"
)

var (
	adminHandler = &http.AdminsHandler{}
	userHandler = &http.UserHandler{}
	chatHandler = &http.ChatHandler{}
	settingHandler = &http.SettingHandler{}
	autoMessageHandler = &http.AutoMessageHandler{}
	autoRuleHandler = &http.AutoRuleHandler{}
	systemRuleHandler = &http.SystemRuleHandler{}
	chatSessionHandler = &http.ChatSessionHandler{}
	dashboardHandler = &http.DashboardHandler{}
	transferHandler = &http.TransferHandler{}
)

func registerAdmin() {
	g := Router.Group("/backend")
	g.POST("/login", http.Login)
	authGroup := g.Group("/")
	authGroup.Use(middleware.Authenticate)

	authGroup.GET("/admins", adminHandler.Index)
	authGroup.GET("/admins/:id", adminHandler.Show)


	authGroup.GET("/me", userHandler.Info)
	authGroup.POST("/me/avatar", userHandler.Avatar)
	authGroup.GET("/me/settings" , userHandler.Setting)
	authGroup.PUT("/me/settings" , userHandler.UpdateSetting)
	authGroup.POST("/me/settings/image", userHandler.SettingImage)

	authGroup.DELETE("/ws/chat-user/:id", chatHandler.RemoveUser)
	authGroup.POST("/ws/req-id", chatHandler.GetReqId)
	authGroup.POST("/ws/chat-user", chatHandler.AcceptUser)
	authGroup.GET("/ws/chat-users", chatHandler.ChatUserList)
	authGroup.POST("/ws/read-all", chatHandler.ReadAll)
	authGroup.POST("/ws/image", chatHandler.Image)
	authGroup.GET("/ws/messages", chatHandler.GetHistoryMessage)
	authGroup.GET("/ws/user/:id", chatHandler.GetUserInfo)
	authGroup.GET("/ws/sessions/:uid", chatHandler.GetHistorySession)
	authGroup.POST("/ws/transfer/:id/cancel", chatHandler.ChatCancelTransfer)
	authGroup.POST("/ws/transfer", chatHandler.Transfer)
	authGroup.GET("/ws/transfer/:id/messages", chatHandler.TransferMessages)


	authGroup.GET("/settings", settingHandler.Index)
	authGroup.PUT("/settings/:name", settingHandler.Update)

	authGroup.GET("/auto-messages", autoMessageHandler.Index)
	authGroup.POST("/auto-messages", autoMessageHandler.Store)
	authGroup.PUT("/auto-messages/:id", autoMessageHandler.Update)
	authGroup.DELETE("/auto-messages/:id", autoMessageHandler.Delete)
	authGroup.GET("/auto-messages/:id", autoMessageHandler.Show)
	authGroup.POST("/auto-messages/image", autoMessageHandler.Image)

	authGroup.GET("/system-auto-rules", systemRuleHandler.Index)
	authGroup.PUT("/system-auto-rules", systemRuleHandler.Update)

	authGroup.GET("/options/messages", autoRuleHandler.MessageOptions)
	authGroup.GET("/options/scenes", autoRuleHandler.SceneOptions)
	authGroup.GET("/options/events", autoRuleHandler.EventOptions)

	authGroup.POST("/auto-rules", autoRuleHandler.Store)
	authGroup.PUT("/auto-rules/:id", autoRuleHandler.Update)
	authGroup.GET("/auto-rules", autoRuleHandler.Index)
	authGroup.GET("/auto-rules/:id", autoRuleHandler.Show)
	authGroup.DELETE("/auto-rules/:id", autoRuleHandler.Delete)

	authGroup.GET("/chat-sessions", chatSessionHandler.Index)
	authGroup.GET("/chat-sessions/:id", chatSessionHandler.Show)
	authGroup.POST("/chat-sessions/:id/cancel", chatSessionHandler.Cancel)

	authGroup.GET("/dashboard/query-info", dashboardHandler.GetUserQueryInfo)
	authGroup.GET("/dashboard/online-info", dashboardHandler.GetOnlineInfo)

	authGroup.GET("/transfers", transferHandler.Index)
	authGroup.POST("/transfers/:id/cancel", transferHandler.Cancel)

	authGroup.GET("/ws", func(c *gin.Context) {
		admin := auth.GetAdmin(c)
		admin.GetSetting()
		conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Log.Error(err)
			util.RespError(c , err.Error())
			return
		}
		client := websocket.NewAdminConn(admin, conn)
		client.Setup()
		websocket.AdminHub.Login(client)
	})
}
