// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CustomerAdminChatSettings is the golang structure of table customer_admin_chat_settings for DAO operations like Where/Data.
type CustomerAdminChatSettings struct {
	g.Meta         `orm:"table:customer_admin_chat_settings, do:true"`
	Id             interface{} //
	AdminId        interface{} //
	Background     interface{} //
	IsAutoAccept   interface{} //
	WelcomeContent interface{} //
	OfflineContent interface{} //
	Name           interface{} //
	LastOnline     *gtime.Time //
	Avatar         interface{} //
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
	DeletedAt      *gtime.Time //
}
