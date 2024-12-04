package service

import (
	"context"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/model"
	"github.com/gorilla/websocket"
)

type (
	IChat interface {
		Accept(ctx context.Context, admin model.CustomerAdmin, sessionId uint) (*api.ChatUser, error)
		Register(ctx context.Context, u any, conn *websocket.Conn) error
		IsOnline(customerId uint, uid uint, t string) bool
		BroadcastWaitingUser(ctx context.Context, customerId uint) error
		UpdateAdminSetting(admin *model.CustomerAdmin)
		RemoveManual(ctx context.Context, uid uint, customerId uint) error
		Transfer(ctx context.Context, fromAdmin *model.CustomerAdmin, toId uint, userId uint, remark string) error
		GetOnlineCount(ctx context.Context, customerId uint) (api.ChatOnlineCount, error)
		GetPlatform(customerId uint, uid uint, t string) string
		GetOnlineAdmin(customerId uint) []api.ChatSimpleUser
		GetOnlineUser(customerId uint) []api.ChatSimpleUser
		NoticeRate(msg *model.CustomerChatMessage)
		NoticeUserRead(customerId uint, uid uint, msgIds []uint)
		NoticeTransfer(ctx context.Context, customer uint, admin uint) error
		NoticeAdminRead(customerId uint, uid uint, msgIds []uint)
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
