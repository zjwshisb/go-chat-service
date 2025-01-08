package grpc

import (
	"context"
	v1 "gf-chat/api/chat/v1"
	"gf-chat/internal/service"
	"github.com/duke-git/lancet/v2/random"
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

func (s *sGrpc) Client(name string) v1.ChatClient {
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
			fn(s.Client(server.GetName()))
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
