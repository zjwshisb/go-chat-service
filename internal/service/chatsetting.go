// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"gf-chat/internal/model"
	"gf-chat/internal/model/entity"
)

type (
	IChatSetting interface {
		GetSubscribeId(customerId uint) string
		First(customerId uint, name string) *entity.CustomerChatSettings
		DefaultAvatarForm(customerId uint) *model.ImageFiled
		GetName(customerId uint) string
		GetAvatar(customerId uint) string
		GetSmsCode(customerId uint) string
		// GetIsAutoTransferManual 是否自动转接人工客服
		GetIsAutoTransferManual(customerId uint) bool
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
