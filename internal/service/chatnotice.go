// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

type (
	IChatNotice interface {
		// Set 设置微信订阅消息id
		Set(uid int) error
		// IsSet 是否设置微信订阅消息id
		IsSet(uid int) bool
		// Remove 移除微信订阅消息id
		Remove(uid int) bool
	}
)

var (
	localChatNotice IChatNotice
)

func ChatNotice() IChatNotice {
	if localChatNotice == nil {
		panic("implement not found for interface IChatNotice, forgot register?")
	}
	return localChatNotice
}

func RegisterChatNotice(i IChatNotice) {
	localChatNotice = i
}
