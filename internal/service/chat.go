package service

import (
	"context"
	api "gf-chat/api/backend/v1"
	"gf-chat/internal/model"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gorilla/websocket"
)

type (
	IChat interface {
		Accept(ctx context.Context, admin model.CustomerAdmin, sessionId uint) (*api.ChatUser, error)
		Register(ctx context.Context, u any, conn *websocket.Conn, platform string) error
		NoticeRate(msg *model.CustomerChatMessage)
		NoticeRead(ctx context.Context, customerId, uid uint, msgIds []uint, t string, forceLocal ...bool) error
		NoticeTransfer(ctx context.Context, customer uint, admin uint, forceLocal ...bool) error
		NoticeUserOnline(ctx context.Context, uid uint, platform string, forceLocal ...bool) error
		NoticeUserOffline(ctx context.Context, uid uint, forceLocal ...bool) error
		NoticeRepeatConnect(ctx context.Context, uid, customerId uint, newUuid string, t string, forceLocal ...bool) error
		BroadcastWaitingUser(ctx context.Context, customerId uint, forceLocal ...bool) error
		BroadcastOnlineAdmins(ctx context.Context, customerId uint, forceLocal ...bool) error
		BroadcastQueueLocation(ctx context.Context, customerId uint, forceLocal ...bool) error
		GetOnlineAdmins(ctx context.Context, customerId uint) ([]api.ChatSimpleUser, error)
		GetOnlineUsers(ctx context.Context, customerId uint) ([]api.ChatSimpleUser, error)
		GetWaitingUsers(ctx context.Context, customerId uint) ([]api.ChatSimpleUser, error)
		GetOnlineUserIds(ctx context.Context, customerId uint, types string, forceLocal ...bool) ([]uint, error)
		GetConnInfo(ctx context.Context, customerId, uid uint, t string, forceLocal ...bool) (exist bool, platform string, err error)
		UpdateAdminSetting(ctx context.Context, id uint, forceLocal ...bool) error
		RemoveManual(ctx context.Context, uid uint, customerId uint) error
		DeliveryMessage(ctx context.Context, msgId uint, types string, forceLocal ...bool) error
		RemoveUser(ctx gctx.Ctx, adminId uint, uid uint) (err error)
		GetUserLastChatTime(ctx gctx.Ctx, adminId uint, uid uint) (uint, error)
		GetUserLimitTime(ctx gctx.Ctx, adminId uint, uid uint) (int64, error)
		UpdateUserLimitTime(ctx gctx.Ctx, adminId uint, uid uint, duration int64) error
		GetInvalidUsers(ctx gctx.Ctx, adminId uint) ([]uint, error)
		GetUsersWithLimitTime(ctx gctx.Ctx, adminId uint) (uids []uint, times []int64, err error)
		IsUserValid(ctx gctx.Ctx, adminId uint, uid uint) (bool, error)
		GetActiveUserCount(ctx gctx.Ctx, adminId uint) (uint, error)
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
