package autorule

import (
	"database/sql"
	"gf-chat/internal/consts"
	"gf-chat/internal/model/do"
	"gf-chat/internal/trait"
	"strings"

	"github.com/gogf/gf/v2/os/gctx"

	"gf-chat/internal/dao"
	"gf-chat/internal/model"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"
)

func init() {
	service.RegisterAutoRule(&sAutoRule{
		Curd: trait.Curd[model.CustomerChatAutoRule]{
			Dao: &dao.CustomerChatAutoRules,
		},
	})
}

type sAutoRule struct {
	trait.Curd[model.CustomerChatAutoRule]
}

func (s *sAutoRule) GetActiveByCustomer(customerId uint) (items []*model.CustomerChatAutoRule) {
	dao.CustomerChatAutoRules.Ctx(gctx.New()).Where(
		do.CustomerChatAutoRules{
			CustomerId: customerId,
			IsSystem:   0,
			IsOpen:     1,
		}).Order("sort").
		WithAll().
		Scan(&items)
	return
}

func (s *sAutoRule) Increment(rule *model.CustomerChatAutoRule) error {
	_, err := dao.CustomerChatAutoRules.Ctx(gctx.New()).WherePri(rule.Id).Increment("count", 1)
	return err
}

func (s *sAutoRule) GetMessage(rule *model.CustomerChatAutoRule) *entity.CustomerChatAutoMessages {
	if rule.MessageId > 0 {
		message := &entity.CustomerChatAutoMessages{}
		err := dao.CustomerChatAutoMessages.Ctx(gctx.New()).Where("id", rule.MessageId).
			Where("customer_id", rule.CustomerId).Scan(message)
		if err == sql.ErrNoRows {
			return nil
		}
		return message
	}
	return nil
}

func (s *sAutoRule) sceneInclude(scenes []*entity.CustomerChatAutoRuleScenes, match string) bool {
	for _, item := range scenes {
		if item.Name == match {
			return true
		}
	}
	return false
}

func (s *sAutoRule) IsMatch(rule *model.CustomerChatAutoRule, scene string, message string) bool {
	switch rule.MatchType {
	case consts.AutoRuleMatchTypeAll:
		if rule.Match == message {
			return s.sceneInclude(rule.Scenes, scene)
		}
	case consts.AutoRuleMatchTypePart:
		if strings.Contains(message, rule.Match) {
			return s.sceneInclude(rule.Scenes, scene)
		}
	}
	return false
}

func (s *sAutoRule) GetEnterRule(customerId uint) (*model.CustomerChatAutoRule, error) {
	return s.GetSystemOne(customerId, consts.AutoRuleMatchEnter)
}

func (s *sAutoRule) GetEnterRuleMessage(customerId uint) (*entity.CustomerChatAutoMessages, error) {
	rule, err := s.GetEnterRule(customerId)
	if err != nil {
		return nil, err
	}
	return s.GetMessage(rule), nil
}

// GetSystemOne 获取系统规则
func (s *sAutoRule) GetSystemOne(customerId uint, match string) (rule *model.CustomerChatAutoRule, err error) {
	err = dao.CustomerChatAutoRules.Ctx(gctx.New()).Where(do.CustomerChatAutoRules{
		CustomerId: customerId,
		Match:      match,
		IsSystem:   1,
	}).Scan(&rule)
	if err != nil {
		return
	}
	if rule == nil {
		err = sql.ErrNoRows
	}
	return
}
