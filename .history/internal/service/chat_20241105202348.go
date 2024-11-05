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
	UpdateAdminSetting(customerId uint, setting *entity.CustomerAdminChatSettings)
	NoticeTransfer(customer, admin uint)
	Accept(admin entity.CustomerAdmins, sessionId uint) (*chat.User, error)
	Register(ctx context.Context, u any, conn *ghttp.WebSocket) error
	IsOnline(customerId uint, uid uint, t string) bool
	BroadcastWaitingUser(customerId uint)
	GetOnlineCount(customerId uint) chat.OnlineCount
	GetPlatform(customerId, uid uint, t string) string
	NoticeRate(msg *entity.CustomerChatMessages)
	NoticeUserRead(customerId, uid uint, msgIds []uint)
	NoticeAdminRead(customerId, uid uint, msgIds []uint)
	Transfer(fromAdmin *entity.CustomerAdmins, toId uint, userId uint, remark string) error
	GetOnlineAdmin(customerId uint) []chat.SimpleUser
	GetOnlineUser(customerId uint) []chat.SimpleUser
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
