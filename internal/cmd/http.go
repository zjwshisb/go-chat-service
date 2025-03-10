package cmd

import (
	"context"
	_ "gf-chat/internal/cache"
	"gf-chat/internal/controller"
	"gf-chat/internal/controller/middleware"
	_ "gf-chat/internal/controller/rule"
	"gf-chat/internal/cron"
	"gf-chat/internal/service"

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
			s.BindHookHandler("/*", ghttp.HookBeforeServe, middleware.Cors)
			controller.RegisterRouter(s)
			go func() {
				cron.Run()
			}()
			if service.Grpc().IsOpen() {
				go func() {
					service.Grpc().StartServer()
				}()
			}
			s.Run()

			return nil
		},
	}
)
