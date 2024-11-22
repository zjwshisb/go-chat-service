// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"gf-chat/internal/model"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type (
	IAdminCtx interface {
		// Init 初始化上下文对象指针到上下文对象中，以便后续的请求流程中可以修改。
		Init(r *ghttp.Request, customCtx *model.AdminCtx)
		// Get 获得上下文变量，如果没有设置，那么返回nil
		Get(ctx context.Context) *model.AdminCtx
		// GetCustomerId 获取客户id
		GetCustomerId(ctx context.Context) uint
		// GetId 获取admin实体
		GetId(ctx context.Context) uint
		// GetUser 获取admin实体
		GetUser(ctx context.Context) *model.CustomerAdmin
		// SetUser 将上下文信息设置到上下文请求中
		SetUser(ctx context.Context, Admin *model.CustomerAdmin)
		// SetData 将上下文信息设置到上下文请求中
		SetData(ctx context.Context, data g.Map)
	}
)

var (
	localAdminCtx IAdminCtx
)

func AdminCtx() IAdminCtx {
	if localAdminCtx == nil {
		panic("implement not found for interface IAdminCtx, forgot register?")
	}
	return localAdminCtx
}

func RegisterAdminCtx(i IAdminCtx) {
	localAdminCtx = i
}
