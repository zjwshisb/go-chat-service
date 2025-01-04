package cmd

import (
	"context"
	_ "gf-chat/internal/cache"
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
			_, err = g.Redis().Do(ctx, "ping")
			if err != nil {
				panic(err)
			}
			s := g.Server()
			g.Dump(1)
			controller.RegisterRouter(s)
			g.Dump(2)
			go func() {
				cron.Run()
			}()
			s.Run()
			g.Dump(3)

			return nil
		},
	}
)
