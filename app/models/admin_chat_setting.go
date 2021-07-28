package models

import "time"

type AdminChatSetting struct {
	Id int64 `json:"id"`
	AdminId int64 `json:"-" gorm:"admin_id;index"`
	Background string `json:"background" gorm:"max:max:512"`
	IsAutoAccept bool `json:"is_auto_accept"`
	WelcomeContent string `json:"welcome_content" gorm:"max:512"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}
