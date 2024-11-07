// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
)

type (
	IAdmin interface {
		All(ctx g.Ctx, w do.CustomerAdmins, with ...any) (items []*model.CustomerAdmin, err error)
		Paginate(ctx context.Context, where *do.CustomerAdmins, p model.QueryInput) (items []*model.CustomerAdmin, total int)
		IsValid(admin *model.CustomerAdmin) error
		GetSetting(ctx context.Context, admin *model.CustomerAdmin) (*entity.CustomerAdminChatSettings, error)
		GetAvatar(model *model.CustomerAdmin) (string, error)
		GetChatName(ctx context.Context, model *model.CustomerAdmin) (string, error)
		First(ctx context.Context, where do.CustomerAdmins) (admin *model.CustomerAdmin, err error)
		GetWechat(adminId uint) *entity.CustomerAdminWechat
		GetDetail(ctx context.Context, id any, month string) ([]*model.ChartLine, *model.AdminDetailInfo, error)
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
