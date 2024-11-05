// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package service

import (
	"context"
	"gf-chat/internal/model"
	"gf-chat/internal/model/chat"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/model/relation"
)

type IChatTransfer interface {
	Paginate(ctx context.Context, w *do.CustomerChatTransfers, p model.QueryInput) (res []*relation.CustomerChatTransfer, total uint)
	FirstEntity(w any) *entity.CustomerChatTransfers
	FirstRelation(w any) *relation.CustomerChatTransfer
	GetRelations(w any) []*relation.CustomerChatTransfer
	RelationToChat(relation *relation.CustomerChatTransfer) chat.Transfer
	Cancel(transfer *relation.CustomerChatTransfer) error
	Accept(transfer *entity.CustomerChatTransfers) error
	Create(fromAdminId, toId, uid uint, remark string) error
	GetUserTransferId(customerId, uid uint) uint
}

var localChatTransfer IChatTransfer

func ChatTransfer() IChatTransfer {
	if localChatTransfer == nil {
		panic("implement not found for interface IChatTransfer, forgot register?")
	}
	return localChatTransfer
}

func RegisterChatTransfer(i IChatTransfer) {
	localChatTransfer = i
}
