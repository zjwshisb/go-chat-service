package backend

import (
	"context"
	baseApi "gf-chat/api"
	"gf-chat/api/v1/backend/user"
	"gf-chat/internal/dao"
	"gf-chat/internal/service"
)

var CMe = &cMe{}

type cMe struct {
}

func (c *cMe) Index(ctx context.Context, req *user.InfoReq) (res *user.InfoRes, err error) {
	admin := service.AdminCtx().GetAdmin(ctx)
	res = &backend.MeRes{
		Id:         admin.Id,
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
		Username:   admin.Username,
	}
	return
}

func (c *cMe) UpdateSetting(ctx context.Context, req *user.UpdateSettingReq) (res *baseApi.NilRes, err error) {
	setting := service.Admin().GetSetting(service.AdminCtx().GetAdmin(ctx).Id)
	setting.Name = req.Name
	setting.Background = req.Background.Path
	setting.Avatar = req.Avatar.Path
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

func (c *cMe) GetSetting(ctx context.Context, req *user.SettingReq) (res *yser.SettingRes, err error) {
	admin := service.AdminCtx().GetAdmin(ctx)
	setting := service.Admin().GetSetting(admin.Id)
	avatar := service.Qiniu().Form(setting.Avatar)
	if avatar == nil {
		avatar = service.ChatSetting().DefaultAvatarForm(admin.CustomerId)
	}
	return &user.SettingRes{
		Background:     service.Qiniu().Form(setting.Background),
		IsAutoAccept:   setting.IsAutoAccept == 1,
		WelcomeContent: setting.WelcomeContent,
		OfflineContent: setting.OfflineContent,
		Name:           setting.Name,
		Avatar:         avatar,
	}, nil
}
