// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CustomerChatMessagesDao is the data access object for table customer_chat_messages.
type CustomerChatMessagesDao struct {
	table   string                      // table is the underlying table name of the DAO.
	group   string                      // group is the database configuration group name of current DAO.
	columns CustomerChatMessagesColumns // columns contains all the column names of Table for convenient usage.
}

// CustomerChatMessagesColumns defines and stores column names for table customer_chat_messages.
type CustomerChatMessagesColumns struct {
	Id         string //
	UserId     string //
	AdminId    string //
	CustomerId string //
	Type       string //
	Content    string //
	ReceivedAt string //
	SendAt     string //
	Source     string //
	SessionId  string //
	ReqId      string //
	ReadAt     string //
	CreatedAt  string //
	UpdatedAt  string //
	DeletedAt  string //
}

// customerChatMessagesColumns holds the columns for table customer_chat_messages.
var customerChatMessagesColumns = CustomerChatMessagesColumns{
	Id:         "id",
	UserId:     "user_id",
	AdminId:    "admin_id",
	CustomerId: "customer_id",
	Type:       "type",
	Content:    "content",
	ReceivedAt: "received_at",
	SendAt:     "send_at",
	Source:     "source",
	SessionId:  "session_id",
	ReqId:      "req_id",
	ReadAt:     "read_at",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
	DeletedAt:  "deleted_at",
}

// NewCustomerChatMessagesDao creates and returns a new DAO object for table data access.
func NewCustomerChatMessagesDao() *CustomerChatMessagesDao {
	return &CustomerChatMessagesDao{
		group:   "default",
		table:   "customer_chat_messages",
		columns: customerChatMessagesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CustomerChatMessagesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *CustomerChatMessagesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *CustomerChatMessagesDao) Columns() CustomerChatMessagesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *CustomerChatMessagesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CustomerChatMessagesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CustomerChatMessagesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
