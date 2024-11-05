// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"gf-chat/internal/dao/internal"
)

// internalCustomerChatAutoMessagesDao is internal type for wrapping internal DAO implements.
type internalCustomerChatAutoMessagesDao = *internal.CustomerChatAutoMessagesDao

// customerChatAutoMessagesDao is the data access object for table customer_chat_auto_messages.
// You can define custom methods on it to extend its functionality as you wish.
type customerChatAutoMessagesDao struct {
	internalCustomerChatAutoMessagesDao
}

var (
	// CustomerChatAutoMessages is globally public accessible object for table customer_chat_auto_messages operations.
	CustomerChatAutoMessages = customerChatAutoMessagesDao{
		internal.NewCustomerChatAutoMessagesDao(),
	}
)

// Fill with you ideas below.
