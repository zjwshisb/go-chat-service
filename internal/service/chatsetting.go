package service

import (
	"context"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/trait"
)

type (
	IChatSetting interface {
		trait.ICurd[entity.CustomerChatSettings]
		DefaultAvatarForm(ctx context.Context, customerId uint) (file *api.File, error error)
		GetName(ctx context.Context, customerId uint) (name string, err error)
		GetAvatar(ctx context.Context, customerId uint) (name string, err error)
		// GetIsAutoTransferManual 是否自动转接人工客服
		GetIsAutoTransferManual(ctx context.Context, customerId uint) (b bool, err error)
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