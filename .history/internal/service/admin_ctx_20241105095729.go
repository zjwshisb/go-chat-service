// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package service

import (
	"context"
	"gf-chat/internal/model"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/model/relation"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type IAdminCtx interface {
	Init(r *ghttp.Request, customCtx *model.AdminCtx)
	Get(ctx context.Context) *model.AdminCtx
	GetKitchens(ctx context.Context) []*relation.CustomerKitchen
	GetSchools(ctx context.Context) []*entity.Schools
	GetCustomerId(ctx context.Context) int
	GetAdmin(ctx context.Context) *entity.CustomerAdmins
	SetUser(ctx context.Context, Admin *entity.CustomerAdmins)
	SetData(ctx context.Context, data g.Map)
}

var localAdminCtx IAdminCtx

func AdminCtx() IAdminCtx {
	if localAdminCtx == nil {
		panic("implement not found for interface IAdminCtx, forgot register?")
	}
	return localAdminCtx
}

func RegisterAdminCtx(i IAdminCtx) {
	localAdminCtx = i
}
