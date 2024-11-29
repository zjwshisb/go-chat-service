// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CustomerAdminChatSettings is the golang structure for table customer_admin_chat_settings.
type CustomerAdminChatSettings struct {
	Id             uint        `json:"ID"              orm:"id"              ` //
	AdminId        uint        `json:"ADMIN_ID"        orm:"admin_id"        ` //
	Background     uint        `json:"BACKGROUND"      orm:"background"      ` //
	IsAutoAccept   uint        `json:"IS_AUTO_ACCEPT"  orm:"is_auto_accept"  ` //
	WelcomeContent string      `json:"WELCOME_CONTENT" orm:"welcome_content" ` //
	OfflineContent string      `json:"OFFLINE_CONTENT" orm:"offline_content" ` //
	Name           string      `json:"NAME"            orm:"name"            ` //
	LastOnline     *gtime.Time `json:"LAST_ONLINE"     orm:"last_online"     ` //
	Avatar         uint        `json:"AVATAR"          orm:"avatar"          ` //
	CreatedAt      *gtime.Time `json:"CREATED_AT"      orm:"created_at"      ` //
	UpdatedAt      *gtime.Time `json:"UPDATED_AT"      orm:"updated_at"      ` //
	DeletedAt      *gtime.Time `json:"DELETED_AT"      orm:"deleted_at"      ` //
}
