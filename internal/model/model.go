package model

import (
	"gf-chat/internal/model/entity"
)

type CustomerChatAutoRule struct {
	entity.CustomerChatAutoRules
	Scenes []*entity.CustomerChatAutoRuleScenes `orm:"with:rule_id=id"`
}

type CustomerAdmin struct {
	entity.CustomerAdmins
	Setting *entity.CustomerAdminChatSettings `orm:"with:admin_id=id"`
}

type CustomerChatMessage struct {
	entity.CustomerChatMessages
	Admin *CustomerAdmin `orm:"with:id=admin_id"`
	User  *entity.Users  `orm:"with:id=user_id"`
}

type CustomerChatTransfer struct {
	entity.CustomerChatTransfers
	User      *entity.Users                `orm:"with:id=user_id"`
	FormAdmin *entity.CustomerAdmins       `orm:"with:id=from_admin_id"`
	ToAdmin   *entity.CustomerAdmins       `orm:"with:id=to_admin_id"`
	ToSession *entity.CustomerChatSessions `orm:"with:id=to_session_id"`
}

type CustomerChatSession struct {
	entity.CustomerChatSessions
	User  *entity.Users          `orm:"with:id=user_id"`
	Admin *entity.CustomerAdmins `orm:"with:id=admin_id"`
}
