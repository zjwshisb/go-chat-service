// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CustomerChatSettingsDao is the data access object for table customer_chat_settings.
type CustomerChatSettingsDao struct {
	table   string                      // table is the underlying table name of the DAO.
	group   string                      // group is the database configuration group name of current DAO.
	columns CustomerChatSettingsColumns // columns contains all the column names of Table for convenient usage.
}

// CustomerChatSettingsColumns defines and stores column names for table customer_chat_settings.
type CustomerChatSettingsColumns struct {
	Id          string //
	Name        string //
	Title       string //
	CustomerId  string //
	Value       string //
	Options     string //
	Type        string //
	Description string //
	CreatedAt   string //
	UpdatedAt   string //
}

// customerChatSettingsColumns holds the columns for table customer_chat_settings.
var customerChatSettingsColumns = CustomerChatSettingsColumns{
	Id:          "id",
	Name:        "name",
	Title:       "title",
	CustomerId:  "customer_id",
	Value:       "value",
	Options:     "options",
	Type:        "type",
	Description: "description",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

// NewCustomerChatSettingsDao creates and returns a new DAO object for table data access.
func NewCustomerChatSettingsDao() *CustomerChatSettingsDao {
	return &CustomerChatSettingsDao{
		group:   "default",
		table:   "customer_chat_settings",
		columns: customerChatSettingsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CustomerChatSettingsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *CustomerChatSettingsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *CustomerChatSettingsDao) Columns() CustomerChatSettingsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *CustomerChatSettingsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CustomerChatSettingsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CustomerChatSettingsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
