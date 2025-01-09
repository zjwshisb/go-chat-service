package service

import (
	"context"
	v1 "gf-chat/api/chat/v1"
	"github.com/gogf/gf/v2/net/gsvc"
)

type (
	IGrpc interface {
		GetServerName() string
		IsOpen() bool
		Client(ctx context.Context, name string) v1.ChatClient
		GetServers(ctx context.Context, inputConfig ...gsvc.SearchInput) ([]gsvc.Service, error)
		CallAll(ctx context.Context, fn func(client v1.ChatClient)) error
		RegisterResolver(ctx context.Context)
		StartServer()
	}
)

var (
	localGrpc IGrpc
)

func Grpc() IGrpc {
	if localFile == nil {
		panic("implement not found for interface IGprc, forgot register?")
	}
	return localGrpc
}

func RegisterGrpc(i IGrpc) {
	localGrpc = i
}
