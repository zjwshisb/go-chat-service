package service

import (
	"context"
	"gf-chat/internal/model"
	"gf-chat/internal/trait"
)

type (
	IChatSession interface {
		trait.ICurd[model.CustomerChatSession]
		Cancel(ctx context.Context, session *model.CustomerChatSession) error
		// Close 关闭会话
		Close(ctx context.Context, session *model.CustomerChatSession, isRemoveUser bool, updateTime bool)
		RelationToChat(session *model.CustomerChatSession) model.ChatSession
		GetUnAcceptModel(ctx context.Context, customerId uint) (res []*model.CustomerChatSession, err error)
		ActiveTransferOne(ctx context.Context, uid uint, adminId uint) (*model.CustomerChatSession, error)
		ActiveNormalOne(ctx context.Context, uid uint, adminId uint) (*model.CustomerChatSession, error)
		ActiveOne(ctx context.Context, uid uint, adminId any, t any) (*model.CustomerChatSession, error)
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
