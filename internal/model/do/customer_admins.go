// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CustomerAdmins is the golang structure of table customer_admins for DAO operations like Where/Data.
type CustomerAdmins struct {
	g.Meta     `orm:"table:customer_admins, do:true"`
	Id         interface{} //
	CustomerId interface{} //
	Username   interface{} //
	Password   interface{} //
	CreatedAt  *gtime.Time //
	UpdatedAt  *gtime.Time //
	DeletedAt  *gtime.Time //
}
