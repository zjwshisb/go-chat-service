package relation

import (
	"gf-chat/internal/model/entity"
)

type CustomerChatAutoRules struct {
	*entity.CustomerChatAutoRules
	Scenes []*entity.CustomerChatAutoRuleScenes `orm:"with:rule_id=id"`
}
