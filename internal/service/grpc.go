package service

import (
	"context"
	v1 "gf-chat/api/chat/v1"
	"github.com/gogf/gf/v2/net/gsvc"
)

type (
	IGrpc interface {
		// GetServerName 当前微服务名称
		GetServerName() string
		// IsOpen 是否启动了微服务
		IsOpen() bool
		// Client 获取对应服务grpc客户端
		Client(ctx context.Context, name string) v1.ChatClient
		// GetServers 获取所有服务列表
		GetServers(ctx context.Context, inputConfig ...gsvc.SearchInput) ([]gsvc.Service, error)
		// CallAll 遍历请求所有服务
		CallAll(ctx context.Context, fn func(client v1.ChatClient)) error
		// RegisterResolver 注册服务发现
		RegisterResolver(ctx context.Context)
		// StartServer 启动微服务
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
