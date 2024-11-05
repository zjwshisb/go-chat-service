// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CustomerChatSettings is the golang structure for table customer_chat_settings.
type CustomerChatSettings struct {
	Id          uint        `json:"ID"          orm:"id"          ` //
	Name        string      `json:"NAME"        orm:"name"        ` //
	Title       string      `json:"TITLE"       orm:"title"       ` //
	CustomerId  uint        `json:"CUSTOMER_ID" orm:"customer_id" ` //
	Value       string      `json:"VALUE"       orm:"value"       ` //
	Options     string      `json:"OPTIONS"     orm:"options"     ` //
	Type        string      `json:"TYPE"        orm:"type"        ` //
	Description string      `json:"DESCRIPTION" orm:"description" ` //
	CreatedAt   *gtime.Time `json:"CREATED_AT"  orm:"created_at"  ` //
	UpdatedAt   *gtime.Time `json:"UPDATED_AT"  orm:"updated_at"  ` //
}
