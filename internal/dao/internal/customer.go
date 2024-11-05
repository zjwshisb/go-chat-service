// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CustomerDao is the data access object for table customer.
type CustomerDao struct {
	table   string          // table is the underlying table name of the DAO.
	group   string          // group is the database configuration group name of current DAO.
	columns CustomerColumns // columns contains all the column names of Table for convenient usage.
}

// CustomerColumns defines and stores column names for table customer.
type CustomerColumns struct {
	Id        string //
	Name      string //
	CreatedAt string //
	UpdatedAt string //
	DeletedAt string //
}

// customerColumns holds the columns for table customer.
var customerColumns = CustomerColumns{
	Id:        "id",
	Name:      "name",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	DeletedAt: "deleted_at",
}

// NewCustomerDao creates and returns a new DAO object for table data access.
func NewCustomerDao() *CustomerDao {
	return &CustomerDao{
		group:   "default",
		table:   "customer",
		columns: customerColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CustomerDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *CustomerDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *CustomerDao) Columns() CustomerColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *CustomerDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CustomerDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CustomerDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
