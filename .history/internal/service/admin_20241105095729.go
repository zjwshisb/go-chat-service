// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package service

import (
	"context"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/model/relation"

	"github.com/gogf/gf/v2/frame/g"
)

type IAdmin interface {
	GetAdmins(ctx g.Ctx, w any) (items []*entity.CustomerAdmins)
	Paginate(ctx context.Context, where *do.CustomerAdmins, p model.QueryInput) (items []*relation.CustomerAdmins, total int)
	IsValid(admin *entity.CustomerAdmins) error
	EntityToRelation(admin *entity.CustomerAdmins) *relation.CustomerAdmins
	GetSetting(adminId uint) *entity.CustomerAdminChatSettings
	GetAvatar(model *relation.CustomerAdmins) string
	GetChatName(model *entity.CustomerAdmins) string
	First(id int) (admin *entity.CustomerAdmins)
	FirstRelation(id int) *relation.CustomerAdmins
	GetWechat(adminId uint) *entity.CustomerAdminWechat
	GetChatAll(customerId int) []*relation.CustomerAdmins
	GetDetail(ctx context.Context, id any, month string) ([]*model.ChartLine, *model.AdminDetailInfo, error)
}

var localAdmin IAdmin

func Admin() IAdmin {
	if localAdmin == nil {
		panic("implement not found for interface IAdmin, forgot register?")
	}
	return localAdmin
}

func RegisterAdmin(i IAdmin) {
	localAdmin = i
}
