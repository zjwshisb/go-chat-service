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
		// All(ctx g.Ctx, where any, with []any) (items []*model.CustomerAdmin, err error)
		// First(ctx g.Ctx, where any) (model *model.CustomerAdmin, err error)
		// Paginate(ctx g.Ctx, where any, p model.QueryInput, with []any) (paginator *model.Paginator[model.CustomerAdmin], err error)
		IsValid(admin *model.CustomerAdmin) error
		GetSetting(ctx g.Ctx, admin *model.CustomerAdmin) (*entity.CustomerAdminChatSettings, error)
		GetAvatar(model *model.CustomerAdmin) (string, error)
		GetChatName(ctx g.Ctx, model *model.CustomerAdmin) (string, error)
		GetWechat(adminId uint) *entity.CustomerAdminWechat
		GetDetail(ctx g.Ctx, id any, month string) ([]*model.ChartLine, *model.AdminDetailInfo, error)
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
