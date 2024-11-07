// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
)

type (
	ISubscribeMsg interface {
		Send(ctx context.Context, customerId uint, uid uint) error
		First(w do.WeappSubscribeMessages) *entity.WeappSubscribeMessages
		GetEntities(customerId uint) []entity.WeappSubscribeMessages
		CheckChatTmpl(e entity.WeappSubscribeMessages) error
		IsTime(key string) bool
		IsThing(key string) bool
		GetParams(e entity.WeappSubscribeMessages) []string
	}
)

var (
	localSubscribeMsg ISubscribeMsg
)

func SubscribeMsg() ISubscribeMsg {
	if localSubscribeMsg == nil {
		panic("implement not found for interface ISubscribeMsg, forgot register?")
	}
	return localSubscribeMsg
}

func RegisterSubscribeMsg(i ISubscribeMsg) {
	localSubscribeMsg = i
}
