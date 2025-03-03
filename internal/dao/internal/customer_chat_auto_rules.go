// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CustomerChatAutoRulesDao is the data access object for table customer_chat_auto_rules.
type CustomerChatAutoRulesDao struct {
	table   string                       // table is the underlying table name of the DAO.
	group   string                       // group is the database configuration group name of current DAO.
	columns CustomerChatAutoRulesColumns // columns contains all the column names of Table for convenient usage.
}

// CustomerChatAutoRulesColumns defines and stores column names for table customer_chat_auto_rules.
type CustomerChatAutoRulesColumns struct {
	Id         string //
	CustomerId string //
	Name       string //
	Match      string //
	MatchType  string //
	ReplyType  string //
	MessageId  string //
	IsSystem   string //
	Sort       string //
	IsOpen     string //
	Scenes     string //
	Count      string //
	CreatedAt  string //
	UpdatedAt  string //
	DeletedAt  string //
}

// customerChatAutoRulesColumns holds the columns for table customer_chat_auto_rules.
var customerChatAutoRulesColumns = CustomerChatAutoRulesColumns{
	Id:         "id",
	CustomerId: "customer_id",
	Name:       "name",
	Match:      "match",
	MatchType:  "match_type",
	ReplyType:  "reply_type",
	MessageId:  "message_id",
	IsSystem:   "is_system",
	Sort:       "sort",
	IsOpen:     "is_open",
	Scenes:     "scenes",
	Count:      "count",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
	DeletedAt:  "deleted_at",
}

// NewCustomerChatAutoRulesDao creates and returns a new DAO object for table data access.
func NewCustomerChatAutoRulesDao() *CustomerChatAutoRulesDao {
	return &CustomerChatAutoRulesDao{
		group:   "default",
		table:   "customer_chat_auto_rules",
		columns: customerChatAutoRulesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CustomerChatAutoRulesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *CustomerChatAutoRulesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *CustomerChatAutoRulesDao) Columns() CustomerChatAutoRulesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *CustomerChatAutoRulesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CustomerChatAutoRulesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CustomerChatAutoRulesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
