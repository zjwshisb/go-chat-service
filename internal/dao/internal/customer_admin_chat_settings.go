// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CustomerAdminChatSettingsDao is the data access object for table customer_admin_chat_settings.
type CustomerAdminChatSettingsDao struct {
	table   string                           // table is the underlying table name of the DAO.
	group   string                           // group is the database configuration group name of current DAO.
	columns CustomerAdminChatSettingsColumns // columns contains all the column names of Table for convenient usage.
}

// CustomerAdminChatSettingsColumns defines and stores column names for table customer_admin_chat_settings.
type CustomerAdminChatSettingsColumns struct {
	Id             string //
	AdminId        string //
	Background     string //
	IsAutoAccept   string //
	WelcomeContent string //
	OfflineContent string //
	Name           string //
	LastOnline     string //
	Avatar         string //
	CreatedAt      string //
	UpdatedAt      string //
	DeletedAt      string //
}

// customerAdminChatSettingsColumns holds the columns for table customer_admin_chat_settings.
var customerAdminChatSettingsColumns = CustomerAdminChatSettingsColumns{
	Id:             "id",
	AdminId:        "admin_id",
	Background:     "background",
	IsAutoAccept:   "is_auto_accept",
	WelcomeContent: "welcome_content",
	OfflineContent: "offline_content",
	Name:           "name",
	LastOnline:     "last_online",
	Avatar:         "avatar",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
	DeletedAt:      "deleted_at",
}

// NewCustomerAdminChatSettingsDao creates and returns a new DAO object for table data access.
func NewCustomerAdminChatSettingsDao() *CustomerAdminChatSettingsDao {
	return &CustomerAdminChatSettingsDao{
		group:   "default",
		table:   "customer_admin_chat_settings",
		columns: customerAdminChatSettingsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CustomerAdminChatSettingsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *CustomerAdminChatSettingsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *CustomerAdminChatSettingsDao) Columns() CustomerAdminChatSettingsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *CustomerAdminChatSettingsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CustomerAdminChatSettingsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CustomerAdminChatSettingsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
