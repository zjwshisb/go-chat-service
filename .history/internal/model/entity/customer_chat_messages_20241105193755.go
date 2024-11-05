// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CustomerChatMessages is the golang structure for table customer_chat_messages.
type CustomerChatMessages struct {
	Id         uint        `json:"ID"          orm:"id"          ` //
	UserId     uint        `json:"USER_ID"     orm:"user_id"     ` //
	AdminId    uint        `json:"ADMIN_ID"    orm:"admin_id"    ` //
	CustomerId uint        `json:"CUSTOMER_ID" orm:"customer_id" ` //
	Type       string      `json:"TYPE"        orm:"type"        ` //
	Content    string      `json:"CONTENT"     orm:"content"     ` //
	ReceivedAt *gtime.Time `json:"RECEIVED_AT" orm:"received_at" ` //
	SendAt     *gtime.Time `json:"SEND_AT"     orm:"send_at"     ` //
	Source     uint         `json:"SOURCE"      orm:"source"      ` //
	SessionId  uint        `json:"SESSION_ID"  orm:"session_id"  ` //
	ReqId      string      `json:"REQ_ID"      orm:"req_id"      ` //
	ReadAt     *gtime.Time `json:"READ_AT"     orm:"read_at"     ` //
	CreatedAt  *gtime.Time `json:"CREATED_AT"  orm:"created_at"  ` //
	UpdatedAt  *gtime.Time `json:"UPDATED_AT"  orm:"updated_at"  ` //
	DeletedAt  *gtime.Time `json:"DELETED_AT"  orm:"deleted_at"  ` //
}
