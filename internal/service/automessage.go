// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"gf-chat/api/v1/backend/automessage"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
)

type (
	IAutoMessage interface {
		First(ctx context.Context, w any) (msg *entity.CustomerChatAutoMessages, err error)
		Paginate(ctx context.Context, w *do.CustomerChatAutoMessages, p model.QueryInput) (items []*entity.CustomerChatAutoMessages, total int)
		All(ctx context.Context, w do.CustomerChatAutoMessages) (items []*entity.CustomerChatAutoMessages, err error)
		EntityToListItem(i entity.CustomerChatAutoMessages) model.AutoMessageListItem
		Update(ctx context.Context, message *entity.CustomerChatAutoMessages, req *automessage.UpdateReq) (count int64, err error)
		Save(ctx context.Context, req *automessage.StoreReq) (id int64, err error)
		ToChatMessage(auto *entity.CustomerChatAutoMessages) (msg *model.CustomerChatMessage, err error)
	}
)

var (
	localAutoMessage IAutoMessage
)

func AutoMessage() IAutoMessage {
	if localAutoMessage == nil {
		panic("implement not found for interface IAutoMessage, forgot register?")
	}
	return localAutoMessage
}

func RegisterAutoMessage(i IAutoMessage) {
	localAutoMessage = i
}
