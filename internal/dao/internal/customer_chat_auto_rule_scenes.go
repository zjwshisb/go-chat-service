// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CustomerChatAutoRuleScenesDao is the data access object for table customer_chat_auto_rule_scenes.
type CustomerChatAutoRuleScenesDao struct {
	table   string                            // table is the underlying table name of the DAO.
	group   string                            // group is the database configuration group name of current DAO.
	columns CustomerChatAutoRuleScenesColumns // columns contains all the column names of Table for convenient usage.
}

// CustomerChatAutoRuleScenesColumns defines and stores column names for table customer_chat_auto_rule_scenes.
type CustomerChatAutoRuleScenesColumns struct {
	Id        string //
	Name      string //
	RuleId    string //
	UpdatedAt string //
	CreatedAt string //
	DeletedAt string //
}

// customerChatAutoRuleScenesColumns holds the columns for table customer_chat_auto_rule_scenes.
var customerChatAutoRuleScenesColumns = CustomerChatAutoRuleScenesColumns{
	Id:        "id",
	Name:      "name",
	RuleId:    "rule_id",
	UpdatedAt: "updated_at",
	CreatedAt: "created_at",
	DeletedAt: "deleted_at",
}

// NewCustomerChatAutoRuleScenesDao creates and returns a new DAO object for table data access.
func NewCustomerChatAutoRuleScenesDao() *CustomerChatAutoRuleScenesDao {
	return &CustomerChatAutoRuleScenesDao{
		group:   "default",
		table:   "customer_chat_auto_rule_scenes",
		columns: customerChatAutoRuleScenesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CustomerChatAutoRuleScenesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *CustomerChatAutoRuleScenesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *CustomerChatAutoRuleScenesDao) Columns() CustomerChatAutoRuleScenesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *CustomerChatAutoRuleScenesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CustomerChatAutoRuleScenesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CustomerChatAutoRuleScenesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
