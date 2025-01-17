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
		// Accept 接入用户
		Accept(ctx context.Context, admin model.CustomerAdmin, sessionId uint) (*api.ChatUser, error)
		// Register hub注册用户websocket
		Register(ctx context.Context, u any, conn *websocket.Conn, platform string) error
		NoticeRate(msg *model.CustomerChatMessage)
		// NoticeRead 通知消息已读
		NoticeRead(ctx context.Context, customerId, uid uint, msgIds []uint, t string, forceLocal ...bool) error
		// NoticeTransfer 通知用户转接
		NoticeTransfer(ctx context.Context, customer uint, admin uint, forceLocal ...bool) error
		// NoticeUserOnline 用户上线通知
		NoticeUserOnline(ctx context.Context, uid uint, platform string, forceLocal ...bool) error
		// NoticeUserOffline 用户下线通知
		NoticeUserOffline(ctx context.Context, uid uint, forceLocal ...bool) error
		// NoticeRepeatConnect 重复连接通知
		NoticeRepeatConnect(ctx context.Context, uid, customerId uint, newUuid string, t string, forceLocal ...bool) error
		// BroadcastWaitingUser 广播等待用户
		BroadcastWaitingUser(ctx context.Context, customerId uint, forceLocal ...bool) error
		// BroadcastOnlineAdmins 广播在线客服
		BroadcastOnlineAdmins(ctx context.Context, customerId uint, forceLocal ...bool) error
		// BroadcastQueueLocation 广播排队位置
		BroadcastQueueLocation(ctx context.Context, customerId uint, forceLocal ...bool) error
		// GetOnlineAdmins 获取在线客服
		GetOnlineAdmins(ctx context.Context, customerId uint) ([]api.ChatSimpleUser, error)
		// GetOnlineUsers 获取在线用户
		GetOnlineUsers(ctx context.Context, customerId uint) ([]api.ChatSimpleUser, error)
		// GetWaitingUsers 获取等待用户
		GetWaitingUsers(ctx context.Context, customerId uint) ([]api.ChatSimpleUser, error)
		// GetOnlineUserIds 获取在线用户id
		GetOnlineUserIds(ctx context.Context, customerId uint, types string, forceLocal ...bool) ([]uint, error)
		// GetConnInfo 获取用户websocket连接信息
		GetConnInfo(ctx context.Context, customerId, uid uint, t string, forceLocal ...bool) (exist bool, platform string, err error)
		// UpdateAdminSetting 更新客服websocket连接struct保存的设置
		UpdateAdminSetting(ctx context.Context, id uint, forceLocal ...bool) error
		// RemoveManual 从待人工介入的用户列表中移除
		RemoveManual(ctx context.Context, uid uint, customerId uint) error
		// DeliveryMessage 投递消息
		DeliveryMessage(ctx context.Context, msgId uint, types string, forceLocal ...bool) error
		// RemoveUser 客服移除用户
		RemoveUser(ctx gctx.Ctx, adminId uint, uid uint) (err error)
		// GetUserLastChatTime 获取用户最后一次聊天时间
		GetUserLastChatTime(ctx gctx.Ctx, adminId uint, uid uint) (uint, error)
		// GetUserLimitTime 获取客服用户限制时间
		GetUserLimitTime(ctx gctx.Ctx, adminId uint, uid uint) (int64, error)
		// UpdateUserLimitTime 更新客服用户限制时间
		UpdateUserLimitTime(ctx gctx.Ctx, adminId uint, uid uint, duration int64) error
		// GetInvalidUsers 获取客服已失效用户
		GetInvalidUsers(ctx gctx.Ctx, adminId uint) ([]uint, error)
		// GetUsersWithLimitTime 获取客服用户和失效时间
		GetUsersWithLimitTime(ctx gctx.Ctx, adminId uint) (uids []uint, times []int64, err error)
		// IsUserValid 判断用户是否失效
		IsUserValid(ctx gctx.Ctx, adminId uint, uid uint) (bool, error)
		// GetActiveUserCount 获取客服有效用户数
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
