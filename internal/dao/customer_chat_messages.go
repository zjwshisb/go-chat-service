// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"gf-chat/internal/dao/internal"
)

// internalCustomerChatMessagesDao is internal type for wrapping internal DAO implements.
type internalCustomerChatMessagesDao = *internal.CustomerChatMessagesDao

// customerChatMessagesDao is the data access object for table customer_chat_messages.
// You can define custom methods on it to extend its functionality as you wish.
type customerChatMessagesDao struct {
	internalCustomerChatMessagesDao
}

var (
	// CustomerChatMessages is globally public accessible object for table customer_chat_messages operations.
	CustomerChatMessages = customerChatMessagesDao{
		internal.NewCustomerChatMessagesDao(),
	}
)

// Fill with you ideas below.
