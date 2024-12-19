package controller

import (
	"gf-chat/internal/controller/backend"
	"gf-chat/internal/controller/frontend"
	"gf-chat/internal/controller/middleware"
	"github.com/gogf/gf/v2/net/ghttp"
)

func RegisterRouter(s *ghttp.Server) {
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.WriteStatus(200, "hello word")
	})
	s.Group("/api", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.Cors, middleware.HandlerResponse).Group("", func(group *ghttp.RouterGroup) {
			group.Group("/user", func(group *ghttp.RouterGroup) {
				group.Bind(
					frontend.CUser,
				)
				group.Middleware(middleware.UserAuth).Group("/chat", func(group *ghttp.RouterGroup) {
					group.Bind(
						frontend.CWs,
						frontend.CChat,
					)
				})
			})
			group.Group("/backend", func(group *ghttp.RouterGroup) {
				group.Bind(
					backend.CCurrentAdmin.Login,
				)
				group.Middleware(middleware.AdminAuth).Group("/", func(group *ghttp.RouterGroup) {
					group.Bind(
						backend.CDashboard,
						backend.CSession,
						backend.CCustomerAdmin,
						backend.CCurrentAdmin.Index,
						backend.CCurrentAdmin.UpdateSetting,
						backend.CCurrentAdmin.GetSetting,
						backend.CAutoMessage,
						backend.CImage,
						backend.CAutoRule,
						backend.CSystemRule,
						backend.CChatSetting,
						backend.CTransfer,
						backend.COption,
						backend.CWs,
						backend.CChat,
						backend.CChatFile,
					)
				})
			})
		})
	})
}
