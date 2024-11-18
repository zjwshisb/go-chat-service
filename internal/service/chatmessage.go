package service

import (
	"context"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/model"
	"gf-chat/internal/trait"
)

type (
	IChatMessage interface {
		trait.ICurd[model.CustomerChatMessage]
		GenReqId() string
		SaveWithUpdate(ctx context.Context, msg *model.CustomerChatMessage) (err error)
		ToRead(ctx context.Context, where any) (int64, error)
		GetAdminName(ctx context.Context, model model.CustomerChatMessage) (string, error)
		RelationToChat(ctx context.Context, message model.CustomerChatMessage) (api.ChatMessage, error)
		GetAvatar(ctx context.Context, model model.CustomerChatMessage) (string, error)
		GetModels(lastId uint, w any, size uint) []*model.CustomerChatMessage
		NewNotice(session *model.CustomerChatSession, content string) *model.CustomerChatMessage
		NewOffline(admin *model.CustomerAdmin) *model.CustomerChatMessage
		NewWelcome(admin *model.CustomerAdmin) *model.CustomerChatMessage
	}
)

var (
	localChatMessage IChatMessage
)

func ChatMessage() IChatMessage {
	if localChatMessage == nil {
		panic("implement not found for interface IChatMessage, forgot register?")
	}
	return localChatMessage
}

func RegisterChatMessage(i IChatMessage) {
	localChatMessage = i
}
