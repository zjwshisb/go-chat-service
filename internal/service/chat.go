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
		BroadcastWaitingUser(ctx context.Context, customerId uint, forceLocal ...bool) error
		BroadcastOnlineAdmins(ctx context.Context, customerId uint, forceLocal ...bool) error
		BroadcastQueueLocation(ctx context.Context, customerId uint, forceLocal ...bool) error
		UpdateAdminSetting(ctx context.Context, id uint, forceLocal ...bool) error
		RemoveManual(ctx context.Context, uid uint, customerId uint) error
		Transfer(ctx context.Context, fromAdmin *model.CustomerAdmin, toId uint, userId uint, remark string) error
		GetOnlineAdmins(ctx context.Context, customerId uint) ([]api.ChatSimpleUser, error)
		GetOnlineUsers(ctx context.Context, customerId uint) ([]api.ChatSimpleUser, error)
		GetWaitingUsers(ctx context.Context, customerId uint) ([]api.ChatSimpleUser, error)
		NoticeRate(msg *model.CustomerChatMessage)
		NoticeRead(ctx context.Context, customerId, uid uint, msgIds []uint, t string, forceLocal ...bool) error
		NoticeTransfer(ctx context.Context, customer uint, admin uint, forceLocal ...bool) error
		NoticeUserOnline(ctx context.Context, uid uint, platform string, forceLocal ...bool) error
		NoticeUserOffline(ctx context.Context, uid uint, forceLocal ...bool) error
		NoticeRepeatConnect(ctx context.Context, uid, customerId uint, newUuid string, t string, forceLocal ...bool) error
		GetConnInfo(ctx context.Context, customerId, uid uint, t string, forceLocal ...bool) (exist bool, platform string)
		DeliveryMessage(ctx context.Context, msgId uint, types string, forceLocal ...bool) error
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
