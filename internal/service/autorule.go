package service

import (
	"gf-chat/internal/model"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/trait"
)

type (
	IAutoRule interface {
		trait.ICurd[model.CustomerChatAutoRule]
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
