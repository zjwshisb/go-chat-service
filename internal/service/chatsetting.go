package service

import (
	"context"
	"gf-chat/internal/model"
	"gf-chat/internal/trait"
)

type (
	IChatSetting interface {
		trait.ICurd[model.CustomerChatSetting]
		GetName(ctx context.Context, customerId uint) (name string, err error)
		GetAvatar(ctx context.Context, customerId uint) (name string, err error)
		// GetIsAutoTransferManual 是否自动转接人工客服
		GetIsAutoTransferManual(ctx context.Context, customerId uint) (b bool, err error)
		RemoveCache(ctx context.Context, customerId uint) error
		GetIsUserShowQueue(ctx context.Context, customerId uint) (isShow bool, err error)
		GetIsUserShowRead(ctx context.Context, customerId uint) (isShow bool, err error)
		GetAiOpen(ctx context.Context, customerId uint) (b bool, err error)
	}
)

var (
	localChatSetting IChatSetting
)

func ChatSetting() IChatSetting {
	if localChatSetting == nil {
		panic("implement not found for interface IChatSetting, forgot register?")
	}
	return localChatSetting
}

func RegisterChatSetting(i IChatSetting) {
	localChatSetting = i
}
