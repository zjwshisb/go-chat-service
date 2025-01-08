package cmd

import (
	"context"
	"github.com/gogf/gf/v2/os/gcmd"
)

var Rpc = &gcmd.Command{
	Name:        "rpc",
	Brief:       "test rpc",
	Description: "test rpc",
	Func: func(ctx context.Context, parser *gcmd.Parser) error {
		//grpcx.Resolver.Register(etcd.New("127.0.0.1:2379@root:123456"))
		//services, _ := gsvc.GetRegistry().Search(ctx, gsvc.SearchInput{
		//	Prefix:   "",
		//	Name:     "",
		//	Version:  "",
		//	Metadata: nil,
		//})
		//for _, s := range services {
		//	conn := grpcx.Client.MustNewGrpcClientConn(s.GetName())
		//	client := v1.NewChatClient(conn)
		//	res, err := client.GetOnlineUsers(ctx, &v1.GetOnlineRequest{CustomerId: 1})
		//	if err != nil {
		//		g.Log().Error(ctx, err)
		//		return nil
		//	}
		//	g.Log().Debug(ctx, "Response:", res.CustomerId)
		//}
		return nil
	},
}
