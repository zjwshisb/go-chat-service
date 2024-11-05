// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SmsTemplatesDao is the data access object for table sms_templates.
type SmsTemplatesDao struct {
	table   string              // table is the underlying table name of the DAO.
	group   string              // group is the database configuration group name of current DAO.
	columns SmsTemplatesColumns // columns contains all the column names of Table for convenient usage.
}

// SmsTemplatesColumns defines and stores column names for table sms_templates.
type SmsTemplatesColumns struct {
	Id         string //
	CustomerId string //
	Type       string // 短信类型
	Name       string // 模板名称，长度为1~30个字符
	Content    string // 模板内容
	Remark     string // 短信模板申请说明
	Status     string // 状态;0:审核中;1:审核通过;2:审核失败
	Code       string // 模板code
	Reason     string // 审核失败理由
	ApprovedAt string // 通过时间
	RejectedAt string // 拒绝时间
	CreatedAt  string //
	UpdatedAt  string //
}

//  smsTemplatesColumns holds the columns for table sms_templates.
var smsTemplatesColumns = SmsTemplatesColumns{
	Id:         "id",
	CustomerId: "customer_id",
	Type:       "type",
	Name:       "name",
	Content:    "content",
	Remark:     "remark",
	Status:     "status",
	Code:       "code",
	Reason:     "reason",
	ApprovedAt: "approved_at",
	RejectedAt: "rejected_at",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
}

// NewSmsTemplatesDao creates and returns a new DAO object for table data access.
func NewSmsTemplatesDao() *SmsTemplatesDao {
	return &SmsTemplatesDao{
		group:   "default",
		table:   "sms_templates",
		columns: smsTemplatesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *SmsTemplatesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *SmsTemplatesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *SmsTemplatesDao) Columns() SmsTemplatesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *SmsTemplatesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *SmsTemplatesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *SmsTemplatesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx *gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
