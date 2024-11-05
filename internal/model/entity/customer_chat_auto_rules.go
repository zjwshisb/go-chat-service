// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CustomerChatAutoRules is the golang structure for table customer_chat_auto_rules.
type CustomerChatAutoRules struct {
	Id         uint        `json:"ID"          orm:"id"          ` //
	CustomerId uint        `json:"CUSTOMER_ID" orm:"customer_id" ` //
	Name       string      `json:"NAME"        orm:"name"        ` //
	Match      string      `json:"MATCH"       orm:"match"       ` //
	MatchType  string      `json:"MATCH_TYPE"  orm:"match_type"  ` //
	ReplyType  string      `json:"REPLY_TYPE"  orm:"reply_type"  ` //
	MessageId  uint        `json:"MESSAGE_ID"  orm:"message_id"  ` //
	IsSystem   uint        `json:"IS_SYSTEM"   orm:"is_system"   ` //
	Sort       uint        `json:"SORT"        orm:"sort"        ` //
	IsOpen     int         `json:"IS_OPEN"     orm:"is_open"     ` //
	Count      int64       `json:"COUNT"       orm:"count"       ` //
	CreatedAt  *gtime.Time `json:"CREATED_AT"  orm:"created_at"  ` //
	UpdatedAt  *gtime.Time `json:"UPDATED_AT"  orm:"updated_at"  ` //
	DeletedAt  *gtime.Time `json:"DELETED_AT"  orm:"deleted_at"  ` //
}
