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
)

type (
	IChatSession interface {
		Get(ctx context.Context, w any) (res []*model.CustomerChatSession)
		Paginate(ctx context.Context, w any, page model.QueryInput) (res []*model.CustomerChatSession, total int)
		Cancel(session *model.CustomerChatSession) error
		// Close 关闭会话
		Close(session *model.CustomerChatSession, isRemoveUser bool, updateTime bool)
		RelationToChat(session *model.CustomerChatSession) model.ChatSession
		FirstRelation(ctx context.Context, w do.CustomerChatSessions) *model.CustomerChatSession
		First(ctx context.Context, w do.CustomerChatSessions) (item *model.CustomerChatSession, err error)
		SaveEntity(ctx context.Context, model *entity.CustomerChatSessions) *entity.CustomerChatSessions
		Create(ctx context.Context, uid uint, customerId uint, t uint) *entity.CustomerChatSessions
		GetUnAcceptModel(ctx context.Context, customerId uint) (res []*model.CustomerChatSession, err error)
		ActiveTransferOne(ctx context.Context, uid uint, adminId uint) (*model.CustomerChatSession, error)
		ActiveNormalOne(ctx context.Context, uid uint, adminId uint) (*model.CustomerChatSession, error)
		ActiveOne(ctx context.Context, uid uint, adminId any, t any) (*model.CustomerChatSession, error)
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
