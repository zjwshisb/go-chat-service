// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Customers is the golang structure for table customers.
type Customers struct {
	Id        uint        `json:"ID"         orm:"id"         ` //
	Name      string      `json:"NAME"       orm:"name"       ` //
	CreatedAt *gtime.Time `json:"CREATED_AT" orm:"created_at" ` //
	UpdatedAt *gtime.Time `json:"UPDATED_AT" orm:"updated_at" ` //
	DeletedAt *gtime.Time `json:"DELETED_AT" orm:"deleted_at" ` //
}
