package grpc

import (
	"context"
	v1 "gf-chat/api/chat/v1"
	"gf-chat/internal/grpc/chat"
	"gf-chat/internal/service"
	"github.com/duke-git/lancet/v2/random"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/contrib/registry/etcd/v2"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gctx"
	"google.golang.org/grpc"
	"sync"
	"time"
)

var name string

var open bool

func init() {
	name = random.RandString(20)
	service.RegisterGrpc(&sGrpc{})

	config, _ := g.Config().Get(gctx.New(), "grpc.open", false)
	open = config.Bool()
}

type sGrpc struct {
}

func (s *sGrpc) clientTimeout(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption,
) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err := invoker(ctx, method, req, reply, cc, opts...)
	return err
}

func (s *sGrpc) StartServer() {
	s.RegisterResolver(gctx.New())
	config := grpcx.Server.NewConfig()
	config.Options = append(config.Options, []grpc.ServerOption{
		grpcx.Server.ChainUnary(
			grpcx.Server.UnaryError,
		)}...,
	)
	config.Name = service.Grpc().GetServerName()
	server := grpcx.Server.New(config)
	chat.Register(server)
	server.Run()
}

func (s *sGrpc) RegisterResolver(ctx context.Context) {
	etcdConfig, err := g.Config().Get(ctx, "etcd.host")
	if err != nil {
		panic(err)
	}
	grpcx.Resolver.Register(etcd.New(etcdConfig.String()))
}

func (s *sGrpc) Client(ctx context.Context, name string) v1.ChatClient {
	servers, err := s.GetServers(ctx)
	if err != nil {
		g.Log().Errorf(ctx, "%+v", err)
		return nil
	}
	_, ok := slice.FindBy(servers, func(index int, item gsvc.Service) bool {
		return item.GetName() == name
	})
	if !ok {
		return nil
	}
	var conn = grpcx.Client.MustNewGrpcClientConn(name, grpcx.Client.ChainUnary(
		s.clientTimeout,
	))
	return v1.NewChatClient(conn)
}

func (s *sGrpc) CallAll(ctx context.Context, fn func(client v1.ChatClient)) error {
	server, err := s.GetServers(ctx)
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	for _, server := range server {
		wg.Add(1)
		go func() {
			fn(s.Client(ctx, server.GetName()))
			wg.Done()
		}()
	}
	wg.Wait()
	return nil
}

func (s *sGrpc) GetServers(ctx context.Context, inputConfig ...gsvc.SearchInput) ([]gsvc.Service, error) {
	input := gsvc.SearchInput{}
	if len(inputConfig) > 0 {
		input = inputConfig[0]
	}
	servers, err := gsvc.GetRegistry().Search(ctx, input)
	if err != nil {
		return nil, err
	}
	return servers, nil
}

func (s *sGrpc) GetServerName() string {
	return name
}

func (s *sGrpc) IsOpen() bool {
	return open
}
