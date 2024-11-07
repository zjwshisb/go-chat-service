// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"gf-chat/internal/model"
	"gf-chat/internal/model/entity"

	"github.com/gorilla/websocket"
)

type (
	IChat interface {
		UpdateAdminSetting(customerId uint, setting *entity.CustomerAdminChatSettings)
		NoticeTransfer(customer uint, admin uint)
		Accept(ctx context.Context, admin model.CustomerAdmin, sessionId uint) (*model.ChatUser, error)
		Register(ctx context.Context, u any, conn *websocket.Conn) error
		IsOnline(customerId uint, uid uint, t string) bool
		BroadcastWaitingUser(customerId uint)
		GetOnlineCount(customerId uint) model.ChatOnlineCount
		GetPlatform(customerId uint, uid uint, t string) string
		NoticeRate(msg *entity.CustomerChatMessages)
		NoticeUserRead(customerId uint, uid uint, msgIds []uint)
		NoticeAdminRead(customerId uint, uid uint, msgIds []uint)
		Transfer(fromAdmin *model.CustomerAdmin, toId uint, userId uint, remark string) error
		GetOnlineAdmin(customerId uint) []model.ChatSimpleUser
		GetOnlineUser(customerId uint) []model.ChatSimpleUser
	}
)

var (
	localChat IChat
)

func Chat() IChat {
	if localChat == nil {
		panic("implement not found for interface IChat, forgot register?")
	}
	return localChat
}

func RegisterChat(i IChat) {
	localChat = i
}
