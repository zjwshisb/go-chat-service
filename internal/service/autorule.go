package service

import (
	"context"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/trait"
)

type (
	IAutoRule interface {
		trait.ICurd[model.CustomerChatAutoRule]
		AllActive(ctx context.Context, customerId uint) (items []*model.CustomerChatAutoRule, err error)
		IncrTriggerCount(ctx context.Context, rule *model.CustomerChatAutoRule) error
		GetMessage(ctx context.Context, rule *model.CustomerChatAutoRule) (msg *model.CustomerChatAutoMessage, err error)
		GetEnterRule(ctx context.Context, customerId uint) (*model.CustomerChatAutoRule, error)
		GetSystemOne(ctx context.Context, customerId uint, match string) (rule *model.CustomerChatAutoRule, err error)
		IsMatch(rule *model.CustomerChatAutoRule, scene string, message string) bool
		Form2Do(form api.AutoRuleForm) *do.CustomerChatAutoRules
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
