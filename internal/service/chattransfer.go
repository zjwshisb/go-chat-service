package service

import (
	"context"
	"gf-chat/internal/model"
	"gf-chat/internal/trait"
)

type (
	IChatTransfer interface {
		trait.ICurd[model.CustomerChatTransfer]
		ToChatTransfer(relation *model.CustomerChatTransfer) model.ChatTransfer
		// Cancel 取消待接入的转接
		Cancel(transfer *model.CustomerChatTransfer) error
		Accept(transfer *model.CustomerChatTransfer) error
		// Create 创建转接
		Create(ctx context.Context, fromAdminId uint, toId uint, uid uint, remark string) error
		GetUserTransferId(customerId uint, uid uint) uint
	}
)

var (
	localChatTransfer IChatTransfer
)

func ChatTransfer() IChatTransfer {
	if localChatTransfer == nil {
		panic("implement not found for interface IChatTransfer, forgot register?")
	}
	return localChatTransfer
}

func RegisterChatTransfer(i IChatTransfer) {
	localChatTransfer = i
}
