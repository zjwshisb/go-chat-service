package service

import (
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/model"
	"gf-chat/internal/trait"

	"context"
	"github.com/gogf/gf/v2/net/ghttp"
)

type (
	IAdmin interface {
		trait.ICurd[model.CustomerAdmin]
		CanAccess(admin *model.CustomerAdmin) bool
		GetAvatar(ctx context.Context, model *model.CustomerAdmin) (string, error)
		GetChatName(ctx context.Context, model *model.CustomerAdmin) (string, error)
		Auth(ctx context.Context, req *ghttp.Request) (admin *model.CustomerAdmin, err error)
		Login(ctx context.Context, request *ghttp.Request) (admin *model.CustomerAdmin, token string, err error)
		UpdateSetting(ctx context.Context, admin *model.CustomerAdmin, form api.CurrentAdminSettingForm) (err error)
		FindSetting(ctx context.Context, adminId uint, withFile bool) (*model.CustomerAdminChatSetting, error)
		GetApiSetting(ctx context.Context, admin *model.CustomerAdmin) (*api.CurrentAdminSetting, error)
		GenSetting(ctx context.Context, admin *model.CustomerAdmin) (*model.CustomerAdminChatSetting, error)
		GetAdminsWithSetting(ctx context.Context, where any) (res []*model.CustomerAdmin, err error)
	}
)

var (
	localAdmin IAdmin
)

func Admin() IAdmin {
	if localAdmin == nil {
		panic("implement not found for interface IAdmin, forgot register?")
	}
	return localAdmin
}

func RegisterAdmin(i IAdmin) {
	localAdmin = i
}
