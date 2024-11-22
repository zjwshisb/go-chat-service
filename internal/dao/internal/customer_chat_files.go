// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CustomerChatFilesDao is the data access object for table customer_chat_files.
type CustomerChatFilesDao struct {
	table   string                   // table is the underlying table name of the DAO.
	group   string                   // group is the database configuration group name of current DAO.
	columns CustomerChatFilesColumns // columns contains all the column names of Table for convenient usage.
}

// CustomerChatFilesColumns defines and stores column names for table customer_chat_files.
type CustomerChatFilesColumns struct {
	Id         string //
	CustomerId string // 客户Id
	Disk       string // 存储引擎
	Path       string // 路径
	Name       string //
	FromModel  string // 来源模型
	FromId     string // 来源id
	Type       string // 文件类型
	CreatedAt  string //
	UpdatedAt  string //
	DeletedAt  string //
}

// customerChatFilesColumns holds the columns for table customer_chat_files.
var customerChatFilesColumns = CustomerChatFilesColumns{
	Id:         "id",
	CustomerId: "customer_id",
	Disk:       "disk",
	Path:       "path",
	Name:       "name",
	FromModel:  "from_model",
	FromId:     "from_id",
	Type:       "type",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
	DeletedAt:  "deleted_at",
}

// NewCustomerChatFilesDao creates and returns a new DAO object for table data access.
func NewCustomerChatFilesDao() *CustomerChatFilesDao {
	return &CustomerChatFilesDao{
		group:   "default",
		table:   "customer_chat_files",
		columns: customerChatFilesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CustomerChatFilesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *CustomerChatFilesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *CustomerChatFilesDao) Columns() CustomerChatFilesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *CustomerChatFilesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CustomerChatFilesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CustomerChatFilesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
