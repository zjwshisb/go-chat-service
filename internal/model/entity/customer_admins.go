// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CustomerAdmins is the golang structure for table customer_admins.
type CustomerAdmins struct {
	Id         uint        `json:"ID"          orm:"id"          ` //
	CustomerId uint        `json:"CUSTOMER_ID" orm:"customer_id" ` //
	Username   string      `json:"USERNAME"    orm:"username"    ` //
	Password   string      `json:"PASSWORD"    orm:"password"    ` //
	CreatedAt  *gtime.Time `json:"CREATED_AT"  orm:"created_at"  ` //
	UpdatedAt  *gtime.Time `json:"UPDATED_AT"  orm:"updated_at"  ` //
	DeletedAt  *gtime.Time `json:"DELETED_AT"  orm:"deleted_at"  ` //
}
