// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CustomerChatTransfersDao is the data access object for table customer_chat_transfers.
type CustomerChatTransfersDao struct {
	table   string                       // table is the underlying table name of the DAO.
	group   string                       // group is the database configuration group name of current DAO.
	columns CustomerChatTransfersColumns // columns contains all the column names of Table for convenient usage.
}

// CustomerChatTransfersColumns defines and stores column names for table customer_chat_transfers.
type CustomerChatTransfersColumns struct {
	Id            string //
	UserId        string //
	FromSessionId string //
	ToSessionId   string //
	FromAdminId   string //
	ToAdminId     string //
	CustomerId    string //
	Remark        string //
	AcceptedAt    string //
	CanceledAt    string //
	CreatedAt     string //
	UpdatedAt     string //
}

// customerChatTransfersColumns holds the columns for table customer_chat_transfers.
var customerChatTransfersColumns = CustomerChatTransfersColumns{
	Id:            "id",
	UserId:        "user_id",
	FromSessionId: "from_session_id",
	ToSessionId:   "to_session_id",
	FromAdminId:   "from_admin_id",
	ToAdminId:     "to_admin_id",
	CustomerId:    "customer_id",
	Remark:        "remark",
	AcceptedAt:    "accepted_at",
	CanceledAt:    "canceled_at",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

// NewCustomerChatTransfersDao creates and returns a new DAO object for table data access.
func NewCustomerChatTransfersDao() *CustomerChatTransfersDao {
	return &CustomerChatTransfersDao{
		group:   "default",
		table:   "customer_chat_transfers",
		columns: customerChatTransfersColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CustomerChatTransfersDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *CustomerChatTransfersDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *CustomerChatTransfersDao) Columns() CustomerChatTransfersColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *CustomerChatTransfersDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CustomerChatTransfersDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CustomerChatTransfersDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
