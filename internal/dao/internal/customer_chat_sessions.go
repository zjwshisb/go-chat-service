// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CustomerChatSessionsDao is the data access object for table customer_chat_sessions.
type CustomerChatSessionsDao struct {
	table   string                      // table is the underlying table name of the DAO.
	group   string                      // group is the database configuration group name of current DAO.
	columns CustomerChatSessionsColumns // columns contains all the column names of Table for convenient usage.
}

// CustomerChatSessionsColumns defines and stores column names for table customer_chat_sessions.
type CustomerChatSessionsColumns struct {
	Id         string //
	UserId     string //
	QueriedAt  string //
	AcceptedAt string //
	CanceledAt string //
	BrokenAt   string //
	CustomerId string //
	AdminId    string //
	Type       string //
	Rate       string //
}

// customerChatSessionsColumns holds the columns for table customer_chat_sessions.
var customerChatSessionsColumns = CustomerChatSessionsColumns{
	Id:         "id",
	UserId:     "user_id",
	QueriedAt:  "queried_at",
	AcceptedAt: "accepted_at",
	CanceledAt: "canceled_at",
	BrokenAt:   "broken_at",
	CustomerId: "customer_id",
	AdminId:    "admin_id",
	Type:       "type",
	Rate:       "rate",
}

// NewCustomerChatSessionsDao creates and returns a new DAO object for table data access.
func NewCustomerChatSessionsDao() *CustomerChatSessionsDao {
	return &CustomerChatSessionsDao{
		group:   "default",
		table:   "customer_chat_sessions",
		columns: customerChatSessionsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CustomerChatSessionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *CustomerChatSessionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *CustomerChatSessionsDao) Columns() CustomerChatSessionsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *CustomerChatSessionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CustomerChatSessionsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CustomerChatSessionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
