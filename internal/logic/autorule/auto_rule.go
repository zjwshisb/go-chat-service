package autorule

import (
	"context"
	"database/sql"
	"gf-chat/internal/consts"
	"gf-chat/internal/model/do"
	"gf-chat/internal/trait"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"strings"

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

func (s *sAutoRule) AllActive(ctx context.Context, customerId uint) ([]*model.CustomerChatAutoRule, error) {
	return s.All(ctx, do.CustomerChatAutoRules{
		CustomerId: customerId,
		IsSystem:   0,
		IsOpen:     1,
	}, g.Slice{model.CustomerChatAutoRule{}.Scenes}, g.Slice{"sort"})
}

func (s *sAutoRule) Increment(ctx context.Context, rule *model.CustomerChatAutoRule) error {
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
	return nil, gerror.New("no message")
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

func (s *sAutoRule) GetEnterRule(ctx context.Context, customerId uint) (*model.CustomerChatAutoRule, error) {
	return s.GetSystemOne(ctx, customerId, consts.AutoRuleMatchEnter)
}

func (s *sAutoRule) GetEnterRuleMessage(ctx context.Context, customerId uint) (*model.CustomerChatAutoMessage, error) {
	rule, err := s.GetEnterRule(ctx, customerId)
	if err != nil {
		return nil, err
	}
	return s.GetMessage(ctx, rule)
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
