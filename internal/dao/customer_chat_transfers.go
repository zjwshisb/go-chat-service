// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"gf-chat/internal/dao/internal"
)

// internalCustomerChatTransfersDao is internal type for wrapping internal DAO implements.
type internalCustomerChatTransfersDao = *internal.CustomerChatTransfersDao

// customerChatTransfersDao is the data access object for table customer_chat_transfers.
// You can define custom methods on it to extend its functionality as you wish.
type customerChatTransfersDao struct {
	internalCustomerChatTransfersDao
}

var (
	// CustomerChatTransfers is globally public accessible object for table customer_chat_transfers operations.
	CustomerChatTransfers = customerChatTransfersDao{
		internal.NewCustomerChatTransfersDao(),
	}
)

// Fill with you ideas below.
