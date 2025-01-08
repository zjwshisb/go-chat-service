package service

import (
	"context"
	api "gf-chat/api/backend/v1"
	"gf-chat/internal/model"

	"github.com/gorilla/websocket"
)

type (
	IChat interface {
		Accept(ctx context.Context, admin model.CustomerAdmin, sessionId uint) (*api.ChatUser, error)
		Register(ctx context.Context, u any, conn *websocket.Conn, platform string) error
		BroadcastWaitingUser(ctx context.Context, customerId uint) error
		UpdateAdminSetting(ctx context.Context, admin *model.CustomerAdmin)
		RemoveManual(ctx context.Context, uid uint, customerId uint) error
		Transfer(ctx context.Context, fromAdmin *model.CustomerAdmin, toId uint, userId uint, remark string) error
		GetOnlineAdmins(ctx context.Context, customerId uint) ([]api.ChatSimpleUser, error)
		GetOnlineUsers(ctx context.Context, customerId uint) ([]api.ChatSimpleUser, error)
		GetWaitingUsers(ctx context.Context, customerId uint) ([]api.ChatSimpleUser, error)
		NoticeRate(msg *model.CustomerChatMessage)
		NoticeRead(ctx context.Context, customerId, uid uint, msgIds []uint, t string, forceLocal ...bool) error
		NoticeTransfer(ctx context.Context, customer uint, admin uint) error
		GetConnInfo(ctx context.Context, customerId, uid uint, t string, forceLocal ...bool) (exist bool, platform string)
		DeliveryAdminMessage(ctx context.Context, msgId uint) error
		DeliveryUserMessage(ctx context.Context, msgId uint) error
		GetOnlineUserIds(ctx context.Context, customerId uint, types string, forceLocal ...bool) ([]uint, error)
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
