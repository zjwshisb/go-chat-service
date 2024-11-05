// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package service

import (
	"context"
	"gf-chat/internal/model"
	"gf-chat/internal/model/chat"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/model/relation"
)

type IChatSession interface {
	Get(ctx context.Context, w any) (res []*relation.CustomerChatSessions)
	Paginate(ctx context.Context, w any, page model.QueryInput) (res []*relation.CustomerChatSessions, total int)
	Cancel(session *entity.CustomerChatSessions) error
	Close(session *entity.CustomerChatSessions, isRemoveUser bool, updateTime bool)
	RelationToChat(model *relation.CustomerChatSessions) chat.Session
	FirstRelation(ctx context.Context, w do.CustomerChatSessions) *relation.CustomerChatSessions
	First(ctx context.Context, w do.CustomerChatSessions) *entity.CustomerChatSessions
	SaveEntity(model *entity.CustomerChatSessions) *entity.CustomerChatSessions
	Create(uid int, customerId int, t int) *entity.CustomerChatSessions
	GetUnAcceptModel(customerId int) (res []*relation.CustomerChatSessions)
	ActiveTransferOne(uid int, adminId int) *entity.CustomerChatSessions
	ActiveNormalOne(uid int, adminId int) *entity.CustomerChatSessions
	ActiveOne(uid int, adminId, t any) *entity.CustomerChatSessions
}

var localChatSession IChatSession

func ChatSession() IChatSession {
	if localChatSession == nil {
		panic("implement not found for interface IChatSession, forgot register?")
	}
	return localChatSession
}

func RegisterChatSession(i IChatSession) {
	localChatSession = i
}
