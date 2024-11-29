// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CustomerChatAutoRuleScenes is the golang structure for table customer_chat_auto_rule_scenes.
type CustomerChatAutoRuleScenes struct {
	Id        uint        `json:"ID"         orm:"id"         ` //
	Name      string      `json:"NAME"       orm:"name"       ` //
	RuleId    uint        `json:"RULE_ID"    orm:"rule_id"    ` //
	UpdatedAt *gtime.Time `json:"UPDATED_AT" orm:"updated_at" ` //
	CreatedAt *gtime.Time `json:"CREATED_AT" orm:"created_at" ` //
	DeletedAt *gtime.Time `json:"DELETED_AT" orm:"deleted_at" ` //
}
