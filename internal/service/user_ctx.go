// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"gf-chat/internal/model"

	"github.com/gogf/gf/v2/net/ghttp"
)

type (
	IUserCtx interface {
		// Init 初始化上下文对象指针到上下文对象中，以便后续的请求流程中可以修改。
		Init(r *ghttp.Request, customCtx *model.UserCtx)
		// Get 获得上下文变量，如果没有设置，那么返回nil
		Get(ctx context.Context) *model.UserCtx
		// GetCustomerId 获取客户id
		GetCustomerId(ctx context.Context) uint
		// GetId 获取客户id
		GetId(ctx context.Context) uint
		// GetUser 获取admin实体
		GetUser(ctx context.Context) *model.User
	}
)

var (
	localUserCtx IUserCtx
)

func UserCtx() IUserCtx {
	if localUserCtx == nil {
		panic("implement not found for interface IUserCtx, forgot register?")
	}
	return localUserCtx
}

func RegisterUserCtx(i IUserCtx) {
	localUserCtx = i
}
