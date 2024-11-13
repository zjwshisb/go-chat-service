package service

import (
	"context"
	"gf-chat/api/v1/backend/automessage"
	"gf-chat/internal/model"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/trait"
)

type (
	IAutoMessage interface {
		trait.ICurd[model.CustomerChatAutoMessage]
		UpdateOne(ctx context.Context, message *model.CustomerChatAutoMessage, req *automessage.UpdateReq) (count int64, err error)
		SaveOne(ctx context.Context, req *automessage.StoreReq) (id int64, err error)
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
