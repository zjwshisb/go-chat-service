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
	IAction interface {
		GetMessage(action *model.ChatAction) (message *model.CustomerChatMessage, err error)
		UnMarshalAction(b []byte) (action *model.ChatAction, err error)
		MarshalAction(action *model.ChatAction) (b []byte, err error)
		String(action *model.ChatAction) string
		NewReceiveAction(msg *model.CustomerChatMessage) *model.ChatAction
		NewReceiptAction(msg *model.CustomerChatMessage) (act *model.ChatAction)
		NewAdminsAction(d any) (act *model.ChatAction)
		NewUserOnline(uid uint, platform string) *model.ChatAction
		NewUserOffline(uid uint) *model.ChatAction
		NewMoreThanOne() *model.ChatAction
		NewOtherLogin() *model.ChatAction
		NewPing() *model.ChatAction
		NewWaitingUsers(i interface{}) *model.ChatAction
		NewWaitingUserCount(count uint) *model.ChatAction
		NewUserTransfer(i interface{}) *model.ChatAction
		NewErrorMessage(msg string) *model.ChatAction
		NewReadAction(msgIds []uint) *model.ChatAction
		NewRateAction(message *entity.CustomerChatMessages) *model.ChatAction
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
