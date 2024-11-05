// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CustomerChatAutoMessages is the golang structure of table customer_chat_auto_messages for DAO operations like Where/Data.
type CustomerChatAutoMessages struct {
	g.Meta     `orm:"table:customer_chat_auto_messages, do:true"`
	Id         interface{} //
	Name       interface{} //
	Type       interface{} //
	Content    interface{} //
	CustomerId interface{} //
	CreatedAt  *gtime.Time //
	UpdatedAt  *gtime.Time //
	DeletedAt  *gtime.Time //
}
