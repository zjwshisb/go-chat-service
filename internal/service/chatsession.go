package service

import (
	"context"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/model"
	"gf-chat/internal/trait"
)

type (
	IChatSession interface {
		trait.ICurd[model.CustomerChatSession]
		Cancel(ctx context.Context, session *model.CustomerChatSession) error
		// Close 关闭会话
		Close(ctx context.Context, session *model.CustomerChatSession, isRemoveUser bool, updateTime bool) error
		RelationToChat(session *model.CustomerChatSession) api.ChatSession
		GetUnAccepts(ctx context.Context, customerId uint) (res []*model.CustomerChatSession, err error)
		FirstTransfer(ctx context.Context, uid uint, adminId uint) (*model.CustomerChatSession, error)
		FirstNormal(ctx context.Context, uid uint, adminId uint) (*model.CustomerChatSession, error)
		FirstActive(ctx context.Context, uid uint, adminId any, t any) (*model.CustomerChatSession, error)
		Create(ctx context.Context, uid uint, customerId uint, t uint) (item *model.CustomerChatSession, err error)
	}
)

var (
	localChatSession IChatSession
)

func ChatSession() IChatSession {
	if localChatSession == nil {
		panic("implement not found for interface IChatSession, forgot register?")
	}
	return localChatSession
}

func RegisterChatSession(i IChatSession) {
	localChatSession = i
}
