// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CustomerChatAutoMessages is the golang structure for table customer_chat_auto_messages.
type CustomerChatAutoMessages struct {
	Id         uint        `json:"ID"          orm:"id"          ` //
	Name       string      `json:"NAME"        orm:"name"        ` //
	Type       string      `json:"TYPE"        orm:"type"        ` //
	Content    string      `json:"CONTENT"     orm:"content"     ` //
	CustomerId uint        `json:"CUSTOMER_ID" orm:"customer_id" ` //
	CreatedAt  *gtime.Time `json:"CREATED_AT"  orm:"created_at"  ` //
	UpdatedAt   *gtime.Time `json:"UPDATE_AT"   orm:"update_at"   ` //
	DeletedAt  *gtime.Time `json:"DELETED_AT"  orm:"deleted_at"  ` //
}
