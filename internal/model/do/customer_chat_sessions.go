// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CustomerChatSessions is the golang structure of table customer_chat_sessions for DAO operations like Where/Data.
type CustomerChatSessions struct {
	g.Meta     `orm:"table:customer_chat_sessions, do:true"`
	Id         interface{} //
	UserId     interface{} //
	QueriedAt  *gtime.Time //
	AcceptedAt *gtime.Time //
	CanceledAt *gtime.Time //
	BrokenAt   *gtime.Time //
	CustomerId interface{} //
	AdminId    interface{} //
	Type       interface{} //
	Rate       interface{} //
}
