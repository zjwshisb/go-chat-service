// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package service

import (
	"context"
	"gf-chat/internal/model/chat"
	"gf-chat/internal/model/entity"

	"github.com/gogf/gf/v2/net/ghttp"
)

type IChat interface {
	UpdateAdminSetting(customerId int, setting *entity.CustomerAdminChatSettings)
	NoticeTransfer(customer, admin int)
	Accept(admin entity.CustomerAdmins, sessionId uint64) (*chat.User, error)
	Register(ctx context.Context, u any, conn *ghttp.WebSocket) error
	IsOnline(customerId int, uid int, t string) bool
	BroadcastWaitingUser(customerId int)
	GetOnlineCount(customerId int) chat.OnlineCount
	GetPlatform(customerId, uid int, t string) string
	NoticeRate(msg *entity.CustomerChatMessages)
	NoticeUserRead(customerId, uid int, msgIds []int64)
	NoticeAdminRead(customerId, uid int, msgIds []int64)
	Transfer(fromAdmin *entity.CustomerAdmins, toId int, userId int, remark string) error
	GetOnlineAdmin(customerId int) []chat.SimpleUser
	GetOnlineUser(customerId int) []chat.SimpleUser
}

var localChat IChat

func Chat() IChat {
	if localChat == nil {
		panic("implement not found for interface IChat, forgot register?")
	}
	return localChat
}

func RegisterChat(i IChat) {
	localChat = i
}
