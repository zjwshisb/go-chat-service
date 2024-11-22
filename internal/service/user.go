package service

import (
	"context"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/trait"

	"github.com/gogf/gf/v2/net/ghttp"
)

type (
	IUser interface {
		trait.ICurd[entity.Users]
		Auth(ctx context.Context, req *ghttp.Request) (*entity.Users, error)
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
