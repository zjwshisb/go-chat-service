// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"gf-chat/internal/model"
	"gf-chat/internal/model/official"
)

type (
	IOfficialMsg interface {
		Chat(admin *model.CustomerAdmin, offline official.Offline) error
		Send(ctx context.Context, message official.Message) error
	}
)

var (
	localOfficialMsg IOfficialMsg
)

func OfficialMsg() IOfficialMsg {
	if localOfficialMsg == nil {
		panic("implement not found for interface IOfficialMsg, forgot register?")
	}
	return localOfficialMsg
}

func RegisterOfficialMsg(i IOfficialMsg) {
	localOfficialMsg = i
}
