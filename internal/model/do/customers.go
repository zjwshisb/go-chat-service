// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Customers is the golang structure of table customers for DAO operations like Where/Data.
type Customers struct {
	g.Meta    `orm:"table:customers, do:true"`
	Id        interface{} //
	Name      interface{} //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
	DeletedAt *gtime.Time //
}
