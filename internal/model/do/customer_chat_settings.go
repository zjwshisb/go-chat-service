// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CustomerChatSettings is the golang structure of table customer_chat_settings for DAO operations like Where/Data.
type CustomerChatSettings struct {
	g.Meta      `orm:"table:customer_chat_settings, do:true"`
	Id          interface{} //
	Name        interface{} //
	Title       interface{} //
	CustomerId  interface{} //
	Value       interface{} //
	Options     interface{} //
	Type        interface{} //
	Description interface{} //
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
}
