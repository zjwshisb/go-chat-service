// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CustomerChatTransfers is the golang structure for table customer_chat_transfers.
type CustomerChatTransfers struct {
	Id            uint        `json:"ID"              orm:"id"              ` //
	UserId        uint        `json:"USER_ID"         orm:"user_id"         ` //
	FromSessionId uint        `json:"FROM_SESSION_ID" orm:"from_session_id" ` //
	ToSessionId   uint        `json:"TO_SESSION_ID"   orm:"to_session_id"   ` //
	FromAdminId   uint        `json:"FROM_ADMIN_ID"   orm:"from_admin_id"   ` //
	ToAdminId     uint        `json:"TO_ADMIN_ID"     orm:"to_admin_id"     ` //
	CustomerId    uint        `json:"CUSTOMER_ID"     orm:"customer_id"     ` //
	Remark        string      `json:"REMARK"          orm:"remark"          ` //
	AcceptedAt    *gtime.Time `json:"ACCEPTED_AT"     orm:"accepted_at"     ` //
	CanceledAt    *gtime.Time `json:"CANCELED_AT"     orm:"canceled_at"     ` //
	CreatedAt     *gtime.Time `json:"CREATED_AT"      orm:"created_at"      ` //
	UpdatedAt     *gtime.Time `json:"UPDATED_AT"      orm:"updated_at"      ` //
}
