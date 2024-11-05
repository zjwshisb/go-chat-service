// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"gf-chat/internal/dao/internal"
)

// internalSmsTemplatesDao is internal type for wrapping internal DAO implements.
type internalSmsTemplatesDao = *internal.SmsTemplatesDao

// smsTemplatesDao is the data access object for table sms_templates.
// You can define custom methods on it to extend its functionality as you wish.
type smsTemplatesDao struct {
	internalSmsTemplatesDao
}

var (
	// SmsTemplates is globally public accessible object for table sms_templates operations.
	SmsTemplates = smsTemplatesDao{
		internal.NewSmsTemplatesDao(),
	}
)

// Fill with you ideas below.
