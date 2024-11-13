// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CustomerChatFiles is the golang structure for table customer_chat_files.
type CustomerChatFiles struct {
	Id         uint        `json:"ID"          orm:"id"          ` //
	CustomerId uint        `json:"CUSTOMER_ID" orm:"customer_id" ` // 客户Id
	Disk       string      `json:"DISK"        orm:"disk"        ` // 存储引擎
	Path       string      `json:"PATH"        orm:"path"        ` // 路径
	FromModel  string      `json:"FROM_MODEL"  orm:"from_model"  ` // 来源模型
	FromId     uint        `json:"FROM_ID"     orm:"from_id"     ` // 来源id
	Type       string      `json:"TYPE"        orm:"type"        ` // 文件类型
	CreatedAt  *gtime.Time `json:"CREATED_AT"  orm:"created_at"  ` //
	UpdatedAt  *gtime.Time `json:"UPDATED_AT"  orm:"updated_at"  ` //
	DeletedAt  *gtime.Time `json:"DELETED_AT"  orm:"deleted_at"  ` //
}
