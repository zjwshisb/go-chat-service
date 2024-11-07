// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
)

type (
	IAutoRule interface {
		Paginate(ctx context.Context, where *do.CustomerChatAutoRules, p model.QueryInput) (items []*entity.CustomerChatAutoRules, total int)
		First(ctx context.Context, w do.CustomerChatAutoRules) (rule *entity.CustomerChatAutoRules, err error)
		GetActiveByCustomer(customerId uint) (items []*model.CustomerChatAutoRule)
		Increment(rule *model.CustomerChatAutoRule) error
		GetMessage(rule *model.CustomerChatAutoRule) *entity.CustomerChatAutoMessages
		IsMatch(rule *model.CustomerChatAutoRule, scene string, message string) bool
		GetEnterRule(customerId uint) (*model.CustomerChatAutoRule, error)
		GetEnterRuleMessage(customerId uint) (*entity.CustomerChatAutoMessages, error)
		// GetSystemOne 获取系统规则
		GetSystemOne(customerId uint, match string) (rule *model.CustomerChatAutoRule, err error)
	}
)

var (
	localAutoRule IAutoRule
)

func AutoRule() IAutoRule {
	if localAutoRule == nil {
		panic("implement not found for interface IAutoRule, forgot register?")
	}
	return localAutoRule
}

func RegisterAutoRule(i IAutoRule) {
	localAutoRule = i
}
