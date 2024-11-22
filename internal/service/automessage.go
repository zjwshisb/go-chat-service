package service

import (
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/model"
	"gf-chat/internal/trait"
)

type (
	IAutoMessage interface {
		trait.ICurd[model.CustomerChatAutoMessage]
		ToChatMessage(auto *model.CustomerChatAutoMessage) (msg *model.CustomerChatMessage, err error)
		Fill(model *model.CustomerChatAutoMessage, form api.AutoMessageForm)
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
