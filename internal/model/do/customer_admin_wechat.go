// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// CustomerAdminWechat is the golang structure of table customer_admin_wechat for DAO operations like Where/Data.
type CustomerAdminWechat struct {
	g.Meta         `orm:"table:customer_admin_wechat, do:true"`
	Id             interface{} //
	AdminId        interface{} //
	OpenId         interface{} //
	OfficialOpenId interface{} //
	Unionid        interface{} //
	Avatar         interface{} //
	CreatedAt      interface{} //
	UpdatedAt      interface{} //
	Info           interface{} //
	IsWechatOnly   interface{} //
}
