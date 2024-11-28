package service

import (
	"context"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/trait"
)

type (
	IAutoMessage interface {
		trait.ICurd[model.CustomerChatAutoMessage]
		ToChatMessage(auto *model.CustomerChatAutoMessage) (msg *model.CustomerChatMessage, err error)
		Form2Do(form api.AutoMessageForm) *do.CustomerChatAutoMessages
		ToApis(ctx context.Context, items []*model.CustomerChatAutoMessage) (resp []*api.AutoMessage, err error)
		ToApi(ctx context.Context, message *model.CustomerChatAutoMessage, files *map[uint]*model.CustomerChatFile) *api.AutoMessage
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
