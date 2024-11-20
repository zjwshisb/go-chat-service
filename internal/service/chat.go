// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/model"
	"gf-chat/internal/model/entity"
	"github.com/gorilla/websocket"
)

type (
	IChat interface {
		UpdateAdminSetting(customerId uint, setting *entity.CustomerAdminChatSettings)
		NoticeTransfer(ctx context.Context, customer uint, admin uint)
		Accept(ctx context.Context, admin model.CustomerAdmin, sessionId uint) (*api.ChatUser, error)
		Register(ctx context.Context, u any, conn *websocket.Conn) error
		IsOnline(customerId uint, uid uint, t string) bool
		BroadcastWaitingUser(ctx context.Context, customerId uint) error
		GetOnlineCount(ctx context.Context, customerId uint) (api.ChatOnlineCount, error)
		GetPlatform(customerId uint, uid uint, t string) string
		NoticeRate(msg *model.CustomerChatMessage)
		NoticeUserRead(customerId uint, uid uint, msgIds []uint)
		NoticeAdminRead(customerId uint, uid uint, msgIds []uint)
		Transfer(ctx context.Context, fromAdmin *model.CustomerAdmin, toId uint, userId uint, remark string) error
		GetOnlineAdmin(customerId uint) []api.ChatSimpleUser
		GetOnlineUser(customerId uint) []api.ChatSimpleUser
		RemoveManual(ctx context.Context, uid uint, customerId uint) error
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
