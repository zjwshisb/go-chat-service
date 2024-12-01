package backend

import (
	"context"
	baseApi "gf-chat/api"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/net/ghttp"
)

var CCurrentAdmin = &cCurrentAdmin{}

type cCurrentAdmin struct {
}

func (c *cCurrentAdmin) Index(ctx context.Context, _ *api.CurrentAdminInfoReq) (res *baseApi.NormalRes[api.CurrentAdminRes], err error) {
	admin := service.AdminCtx().GetUser(ctx)
	res = baseApi.NewResp(api.CurrentAdminRes{
		Id:         admin.Id,
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
		Username:   admin.Username,
	})
	return
}

func (c *cCurrentAdmin) UpdateSetting(ctx context.Context, req *api.CurrentAdminSettingUpdateReq) (res *baseApi.NilRes, err error) {
	setting, err := service.Admin().GetSetting(ctx, service.AdminCtx().GetUser(ctx))
	if err != nil {
		return nil, err
	}
	err = service.Admin().UpdateSetting(ctx, service.AdminCtx().GetUser(ctx), req.CurrentAdminSettingForm)
	if err != nil {
		return
	}
	service.Chat().UpdateAdminSetting(service.AdminCtx().GetCustomerId(ctx), setting)
	return baseApi.NewNilResp(), nil
}

func (c *cCurrentAdmin) GetSetting(ctx context.Context, _ *api.CurrentAdminSettingReq) (res *baseApi.NormalRes[*api.CurrentAdminSetting], err error) {
	setting, err := service.Admin().GetSetting(ctx, service.AdminCtx().GetUser(ctx))
	if err != nil {
		return nil, err
	}
	return baseApi.NewResp(setting), nil
}

func (c *cCurrentAdmin) Login(ctx context.Context, _ *api.LoginReq) (res *baseApi.NormalRes[api.LoginRes], err error) {
	request := ghttp.RequestFromCtx(ctx)
	_, token, err := service.Admin().Login(ctx, request)
	if err != nil {
		return
	}
	return baseApi.NewResp(api.LoginRes{
		Token: token,
	}), nil
}
