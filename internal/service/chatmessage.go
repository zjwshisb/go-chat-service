// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"database/sql"
	"gf-chat/internal/model"
	"gf-chat/internal/trait"
	"gf-chat/internal/model/entity"
)

type (
	IChatMessage interface {
		trait.ICurd[model.CustomerChatMessage]
		GenReqId() string
		SaveWithUpdate(ctx context.Context, msg *model.CustomerChatMessage) (err error)
		EntityToRelation(msg *entity.CustomerChatMessages) *model.CustomerChatMessage
		ChangeToRead(msgId []uint) (sql.Result, error)
		GetAdminName(model model.CustomerChatMessage) string
		RelationToChat(message model.CustomerChatMessage) model.ChatMessage
		GetAvatar(model model.CustomerChatMessage) string
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
