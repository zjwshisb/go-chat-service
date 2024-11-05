// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package service

import (
	"context"
	"gf-chat/internal/model"
	"gf-chat/internal/model/entity"

	"github.com/gogf/gf/v2/net/ghttp"
)

type IUserCtx interface {
	Init(r *ghttp.Request, customCtx *model.UserCtx)
	Get(ctx context.Context) *model.UserCtx
	GetCustomerId(ctx context.Context) int
	GetUserApp(ctx context.Context) *entity.UserApps
	GetUser(ctx context.Context) *entity.Users
}

var localUserCtx IUserCtx

func UserCtx() IUserCtx {
	if localUserCtx == nil {
		panic("implement not found for interface IUserCtx, forgot register?")
	}
	return localUserCtx
}

func RegisterUserCtx(i IUserCtx) {
	localUserCtx = i
}
