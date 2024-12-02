package cmd

import (
	"context"
	"gf-chat/internal/controller"
	_ "gf-chat/internal/controller/rule"
	"gf-chat/internal/cron"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/frame/g"
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
			controller.RegisterRouter(s)
			go func() {
				cron.Run()
			}()
			s.Run()
			return nil
		},
	}
)
