package service

import (
	"context"
	api "gf-chat/api/backend/v1"
	"gf-chat/internal/model"
	"gf-chat/internal/trait"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/gogf/gf/v2/net/ghttp"
)

type (
	IUser interface {
		trait.ICurd[model.User]
		// Auth 用户认证
		Auth(ctx g.Ctx, req *ghttp.Request) (admin *model.User, err error)
		// Login 用户登录
		Login(ctx context.Context, request *ghttp.Request) (admin *model.User, token string, err error)
		// GetInfo 获取用户信息
		GetInfo(ctx context.Context, user *model.User) ([]api.UserInfoItem, error)
		// GetActiveCount 获取某天活跃用户
		GetActiveCount(ctx context.Context, customerId uint, date *gtime.Time) (count int, err error)
	}
)

var (
	localUser IUser
)

func User() IUser {
	if localUser == nil {
		panic("implement not found for interface IUser, forgot register?")
	}
	return localUser
}

func RegisterUser(i IUser) {
	localUser = i
}
