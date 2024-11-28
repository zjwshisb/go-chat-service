package cmd

import (
	"context"
	"gf-chat/internal/controller/backend"
	"gf-chat/internal/controller/frontend"
	"gf-chat/internal/controller/middleware"
	_ "gf-chat/internal/controller/rule"
	"gf-chat/internal/cron"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = &gcmd.Command{
		Name:        "main",
		Brief:       "start http server",
		Description: "this is the command entry for starting your process",
	}
	Http = &gcmd.Command{
		Name:  "http",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			s.BindHandler("/", func(r *ghttp.Request) {
				r.Response.WriteStatus(200, "hello word")
			})
			s.Group("/api", func(group *ghttp.RouterGroup) {
				group.Group("/user", func(group *ghttp.RouterGroup) {
					group.Middleware(middleware.Cors, middleware.UserAuth, middleware.HandlerResponse).Group("/chat", func(group *ghttp.RouterGroup) {
						group.Bind(
							frontend.CWs,
							frontend.CChat,
						)
					})
				})
				group.Group("/backend", func(group *ghttp.RouterGroup) {
					group.Middleware(middleware.Cors, middleware.HandlerResponse).
						Group("/", func(group *ghttp.RouterGroup) {
							group.Bind(
								backend.CUser.Login,
							)
							group.Middleware(middleware.AdminAuth).Group("/", func(group *ghttp.RouterGroup) {
								group.Bind(
									backend.CDashboard,
									backend.CSession,
									backend.CCustomerAdmin,
									backend.CUser.Index,
									backend.CUser.UpdateSetting,
									backend.CUser.GetSetting,
									backend.CAutoMessage,
									backend.CImage,
									backend.CAutoRule,
									backend.CSystemRule,
									backend.CChatSetting,
									backend.CTransfer,
									backend.COption,
									// backend.CWs,
									// backend.CChat,
								)
							})
						})
				})
			})

			go func() {
				cron.Run()
			}()
			s.Run()
			return nil
		},
	}
)
