// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CustomerChatSessions is the golang structure for table customer_chat_sessions.
type CustomerChatSessions struct {
	Id         uint        `json:"ID"          orm:"id"          ` //
	UserId     uint        `json:"USER_ID"     orm:"user_id"     ` //
	QueriedAt  *gtime.Time `json:"QUERIED_AT"  orm:"queried_at"  ` //
	AcceptedAt *gtime.Time `json:"ACCEPTED_AT" orm:"accepted_at" ` //
	CanceledAt *gtime.Time `json:"CANCELED_AT" orm:"canceled_at" ` //
	BrokenAt   *gtime.Time `json:"BROKEN_AT"   orm:"broken_at"   ` //
	CustomerId uint        `json:"CUSTOMER_ID" orm:"customer_id" ` //
	AdminId    uint        `json:"ADMIN_ID"    orm:"admin_id"    ` //
	Type       uint        `json:"TYPE"        orm:"type"        ` //
	Rate       uint        `json:"RATE"        orm:"rate"        ` //
}
