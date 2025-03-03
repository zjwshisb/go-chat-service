package service

import (
	"context"
	api "gf-chat/api/backend/v1"
	"gf-chat/internal/model"
	"gf-chat/internal/trait"
)

type (
	IChatTransfer interface {
		trait.ICurd[model.CustomerChatTransfer]
		ToApi(relation *model.CustomerChatTransfer) api.ChatTransfer
		// Cancel 取消待接入的转接
		Cancel(ctx context.Context, transfer *model.CustomerChatTransfer) error
		Accept(ctx context.Context, transfer *model.CustomerChatTransfer) error
		// Create 创建转接
		Create(ctx context.Context, fromAdmin *model.CustomerAdmin, toId uint, userId uint, remark string) error
		GetUserTransferId(ctx context.Context, customerId uint, uid uint) (uint, error)
		RemoveUser(ctx context.Context, customerId, uid uint) error
		IsInTransfer(ctx context.Context, customerId uint, uid uint) (bool, error)
		FirstActive(ctx context.Context, where any) (transfer *model.CustomerChatTransfer, err error)
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
