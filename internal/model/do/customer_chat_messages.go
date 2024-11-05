// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CustomerChatMessages is the golang structure of table customer_chat_messages for DAO operations like Where/Data.
type CustomerChatMessages struct {
	g.Meta     `orm:"table:customer_chat_messages, do:true"`
	Id         interface{} //
	UserId     interface{} //
	AdminId    interface{} //
	CustomerId interface{} //
	Type       interface{} //
	Content    interface{} //
	ReceivedAt *gtime.Time //
	SendAt     *gtime.Time //
	Source     interface{} //
	SessionId  interface{} //
	ReqId      interface{} //
	ReadAt     *gtime.Time //
	CreatedAt  *gtime.Time //
	UpdatedAt  *gtime.Time //
	DeletedAt  *gtime.Time //
}
