// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/model"
)

type (
	IAction interface {
		GetMessage(action *api.ChatAction) (message *model.CustomerChatMessage, err error)
		UnMarshalAction(b []byte) (action *api.ChatAction, err error)
		MarshalAction(action *api.ChatAction) (b []byte, err error)
		String(action *api.ChatAction) string
		NewReceiveAction(msg *model.CustomerChatMessage) *api.ChatAction
		NewReceiptAction(msg *model.CustomerChatMessage) (act *api.ChatAction)
		NewAdminsAction(d any) (act *api.ChatAction)
		NewUserOnline(uid uint, platform string) *api.ChatAction
		NewUserOffline(uid uint) *api.ChatAction
		NewMoreThanOne() *api.ChatAction
		NewOtherLogin() *api.ChatAction
		NewPing() *api.ChatAction
		NewWaitingUsers(i interface{}) *api.ChatAction
		NewWaitingUserCount(count uint) *api.ChatAction
		NewUserTransfer(i interface{}) *api.ChatAction
		NewErrorMessage(msg string) *api.ChatAction
		NewReadAction(msgIds []uint) *api.ChatAction
		NewRateAction(message *model.CustomerChatMessage) *api.ChatAction
	}
)

var (
	localAction IAction
)

func Action() IAction {
	if localAction == nil {
		panic("implement not found for interface IAction, forgot register?")
	}
	return localAction
}

func RegisterAction(i IAction) {
	localAction = i
}
