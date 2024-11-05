// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CustomerChatAutoRuleScenes is the golang structure of table customer_chat_auto_rule_scenes for DAO operations like Where/Data.
type CustomerChatAutoRuleScenes struct {
	g.Meta    `orm:"table:customer_chat_auto_rule_scenes, do:true"`
	Id        interface{} //
	Name      interface{} //
	RuleId    interface{} //
	UpdatedAt *gtime.Time //
	CreatedAt *gtime.Time //
	DeletedAt *gtime.Time //
}
