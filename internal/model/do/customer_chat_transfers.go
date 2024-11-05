// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CustomerChatTransfers is the golang structure of table customer_chat_transfers for DAO operations like Where/Data.
type CustomerChatTransfers struct {
	g.Meta        `orm:"table:customer_chat_transfers, do:true"`
	Id            interface{} //
	UserId        interface{} //
	FromSessionId interface{} //
	ToSessionId   interface{} //
	FromAdminId   interface{} //
	ToAdminId     interface{} //
	CustomerId    interface{} //
	Remark        interface{} //
	AcceptedAt    *gtime.Time //
	CanceledAt    *gtime.Time //
	CreatedAt     *gtime.Time //
	UpdatedAt     *gtime.Time //
}
