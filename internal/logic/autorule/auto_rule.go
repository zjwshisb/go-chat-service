package autorule

import (
	"context"
	"database/sql"
	"gf-chat/internal/consts"
	"gf-chat/internal/model/do"
	"gf-chat/internal/trait"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"strings"

	"gf-chat/internal/dao"
	"gf-chat/internal/model"
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

func (s *sAutoRule) AllActive(ctx context.Context, customerId uint) ([]*model.CustomerChatAutoRule, error) {
	return s.All(ctx, do.CustomerChatAutoRules{
		CustomerId: customerId,
		IsSystem:   0,
		IsOpen:     1,
	}, g.Slice{model.CustomerChatAutoRule{}.Scenes}, g.Slice{"sort"})
}

func (s *sAutoRule) IncrTriggerCount(ctx context.Context, rule *model.CustomerChatAutoRule) error {
	_, err := dao.CustomerChatAutoRules.Ctx(ctx).WherePri(rule.Id).Increment("count", 1)
	return err
}

func (s *sAutoRule) GetMessage(ctx context.Context, rule *model.CustomerChatAutoRule) (msg *model.CustomerChatAutoMessage, err error) {
	if rule.MessageId > 0 {
		msg, err = service.AutoMessage().Find(ctx, rule.MessageId)
		if err != nil {
			return
		}
		return
	}
	return nil, gerror.NewCode(gcode.CodeNotFound)
}

func (s *sAutoRule) sceneInclude(rule *model.CustomerChatAutoRule, match string) bool {
	for _, item := range rule.Scenes {
		if item == match {
			return true
		}
	}
	return false
}

func (s *sAutoRule) IsMatch(rule *model.CustomerChatAutoRule, scene string, message string) bool {
	switch rule.MatchType {
	case consts.AutoRuleMatchTypeAll:
		if rule.Match == message {
			return s.sceneInclude(rule, scene)
		}
	case consts.AutoRuleMatchTypePart:
		if strings.Contains(message, rule.Match) {
			return s.sceneInclude(rule, scene)
		}
	}
	return false
}

func (s *sAutoRule) GetEnterRule(ctx context.Context, customerId uint) (*model.CustomerChatAutoRule, error) {
	return s.GetSystemOne(ctx, customerId, consts.AutoRuleMatchEnter)
}

// GetSystemOne 获取系统规则
func (s *sAutoRule) GetSystemOne(ctx context.Context, customerId uint, match string) (rule *model.CustomerChatAutoRule, err error) {
	err = dao.CustomerChatAutoRules.Ctx(ctx).Where(do.CustomerChatAutoRules{
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
