package service

import (
	"gf-chat/internal/model"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/trait"

	"github.com/gogf/gf/v2/frame/g"
)

type (
	IAdmin interface {
		trait.ICurd[model.CustomerAdmin]
		CanAccess(admin *model.CustomerAdmin) bool
		GetSetting(ctx g.Ctx, admin *model.CustomerAdmin) (*entity.CustomerAdminChatSettings, error)
		GetAvatar(ctx g.Ctx, model *model.CustomerAdmin) (string, error)
		GetChatName(ctx g.Ctx, model *model.CustomerAdmin) (string, error)
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
