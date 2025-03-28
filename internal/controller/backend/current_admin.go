package backend

import (
	"context"
	baseApi "gf-chat/api"
	api "gf-chat/api/backend/v1"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/net/ghttp"
)

var CCurrentAdmin = &cCurrentAdmin{}

type cCurrentAdmin struct {
}

func (c *cCurrentAdmin) Index(ctx context.Context, _ *api.CurrentAdminInfoReq) (res *baseApi.NormalRes[api.CurrentAdmin], err error) {
	admin := service.AdminCtx().GetUser(ctx)
	res = baseApi.NewResp(api.CurrentAdmin{
		Id:         admin.Id,
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
		Username:   admin.Username,
	})
	return
}

func (c *cCurrentAdmin) UpdateSetting(ctx context.Context, req *api.CurrentAdminSettingUpdateReq) (res *baseApi.NilRes, err error) {
	admin := service.AdminCtx().GetUser(ctx)
	err = service.Admin().UpdateSetting(ctx, admin, req.CurrentAdminSettingForm)
	if err != nil {
		return
	}
	err = service.Chat().UpdateAdminSetting(ctx, admin.Id)
	if err != nil {
		return
	}
	return baseApi.NewNilResp(), nil
}

func (c *cCurrentAdmin) GetSetting(ctx context.Context, _ *api.CurrentAdminSettingReq) (res *baseApi.NormalRes[*api.CurrentAdminSetting], err error) {
	setting, err := service.Admin().GetApiSetting(ctx, service.AdminCtx().GetUser(ctx))
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
