// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CustomerChatFiles is the golang structure of table customer_chat_files for DAO operations like Where/Data.
type CustomerChatFiles struct {
	g.Meta     `orm:"table:customer_chat_files, do:true"`
	Id         interface{} //
	CustomerId interface{} // 客户Id
	Disk       interface{} // 存储引擎
	Path       interface{} // 路径
	FromModel  interface{} // 来源模型
	FromId     interface{} // 来源id
	Type       interface{} // 文件类型
	CreatedAt  *gtime.Time //
	UpdatedAt  *gtime.Time //
	DeletedAt  *gtime.Time //
}
