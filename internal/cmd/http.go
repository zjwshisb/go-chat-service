package cmd

import (
	"context"
	_ "gf-chat/internal/cache"
	"gf-chat/internal/controller"
	_ "gf-chat/internal/controller/rule"
	"gf-chat/internal/cron"
	"gf-chat/internal/grpc/chat"
	"gf-chat/internal/service"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/contrib/registry/etcd/v2"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
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
			controller.RegisterRouter(s)
			go func() {
				cron.Run()
			}()
			if service.Grpc().IsOpen() {
				go func() {
					grpcx.Resolver.Register(etcd.New("127.0.0.1:2379@root:123456"))
					config := grpcx.Server.NewConfig()
					config.Name = service.Grpc().GetServerName()
					s := grpcx.Server.New(config)
					chat.Register(s)
					s.Run()
				}()
			}

			s.Run()

			return nil
		},
	}
)
