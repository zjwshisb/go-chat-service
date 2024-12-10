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
		ToRead(ctx context.Context, where any) (int64, error)
		GetAdminName(ctx context.Context, model *model.CustomerChatMessage) (string, error)
		ToApi(ctx context.Context, message *model.CustomerChatMessage) (*api.ChatMessage, error)
		GetAvatar(ctx context.Context, model *model.CustomerChatMessage) (string, error)
		NewNotice(session *model.CustomerChatSession, content string) *model.CustomerChatMessage
		NewOffline(ctx context.Context, admin *model.CustomerAdmin) (msg *model.CustomerChatMessage, err error)
		NewWelcome(ctx context.Context, admin *model.CustomerAdmin) (msg *model.CustomerChatMessage, err error)
		Insert(ctx context.Context, message *model.CustomerChatMessage) (*model.CustomerChatMessage, error)
		GetLastGroupByUsers(ctx context.Context, adminId uint, uids []uint) (res []*model.CustomerChatMessage, err error)
		GetUnreadCountGroupByUsers(ctx context.Context, uids []uint, w any) (res []model.UnreadCount, err error)
		IsTypeValid(types string) (valid bool)
		IsFileType(types string) (valid bool)
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
