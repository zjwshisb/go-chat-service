// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CustomerChatAutoRules is the golang structure of table customer_chat_auto_rules for DAO operations like Where/Data.
type CustomerChatAutoRules struct {
	g.Meta     `orm:"table:customer_chat_auto_rules, do:true"`
	Id         interface{} //
	CustomerId interface{} //
	Name       interface{} //
	Match      interface{} //
	MatchType  interface{} //
	ReplyType  interface{} //
	MessageId  interface{} //
	IsSystem   interface{} //
	Sort       interface{} //
	IsOpen     interface{} //
	Scenes     interface{} //
	Count      interface{} //
	CreatedAt  *gtime.Time //
	UpdatedAt  *gtime.Time //
	DeletedAt  *gtime.Time //
}
