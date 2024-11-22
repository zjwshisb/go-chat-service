package backend

import (
	"context"
	baseApi "gf-chat/api"
	"gf-chat/api/v1/backend"
	"gf-chat/internal/dao"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/net/ghttp"
)

var CUser = &cUser{}

type cUser struct {
}

func (c *cUser) Index(ctx context.Context, req *backend.CurrentUserInfoReq) (res *baseApi.NormalRes[backend.CurrentUserRes], err error) {
	admin := service.AdminCtx().GetUser(ctx)
	res = baseApi.NewResp(backend.CurrentUserRes{
		Id:         admin.Id,
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
		Username:   admin.Username,
	})
	return
}

func (c *cUser) UpdateSetting(ctx context.Context, req *backend.CurrentUserUpdateSettingReq) (res *baseApi.NilRes, err error) {
	setting, err := service.Admin().GetSetting(ctx, service.AdminCtx().GetUser(ctx))
	if err != nil {
		return nil, err
	}
	setting.Name = req.Name
	//setting.Background = req.Background.Path
	//setting.Avatar = req.Avatar.Path
	if req.IsAutoAccept {
		setting.IsAutoAccept = 1
	} else {
		setting.IsAutoAccept = 0
	}
	setting.WelcomeContent = req.WelcomeContent
	setting.OfflineContent = req.OfflineContent
	dao.CustomerAdminChatSettings.Ctx(ctx).Save(setting)
	service.Chat().UpdateAdminSetting(service.AdminCtx().GetCustomerId(ctx), setting)
	return &baseApi.NilRes{}, nil
}

func (c *cUser) GetSetting(ctx context.Context, req *backend.CurrentUserSettingReq) (res *baseApi.NormalRes[backend.CurrentUserSettingRes], err error) {
	//admin := service.AdminCtx().GetUser(ctx)
	setting, err := service.Admin().GetSetting(ctx, service.AdminCtx().GetUser(ctx))
	if err != nil {
		return nil, err
	}
	//avatar := service.Qiniu().Form(setting.Avatar)
	//if avatar == nil {
	//avatar := service.ChatSetting().DefaultAvatarForm(admin.CustomerId)
	//}
	return baseApi.NewResp(backend.CurrentUserSettingRes{
		//Background:     service.Qiniu().Form(setting.Background),
		IsAutoAccept:   setting.IsAutoAccept == 1,
		WelcomeContent: setting.WelcomeContent,
		OfflineContent: setting.OfflineContent,
		Name:           setting.Name,
		//Avatar:         avatar,
	}), nil
}

func (c *cUser) Login(ctx context.Context, r *backend.LoginReq) (res *baseApi.NormalRes[backend.LoginRes], err error) {
	request := ghttp.RequestFromCtx(ctx)
	_, token, err := service.Admin().Login(ctx, request)
	if err != nil {
		return
	}
	return baseApi.NewResp(backend.LoginRes{
		Token: token,
	}), nil
}
