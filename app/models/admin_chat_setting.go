package models

import (
	"time"
	"ws/app/databases"
)

type AdminChatSetting struct {
	Id int64 `json:"id"`
	AdminId int64 `json:"-" gorm:"uniqueIndex"`
	Background string `json:"background" gorm:"size:512"`
	IsAutoAccept bool `json:"is_auto_accept"`
	WelcomeContent string `json:"welcome_content" gorm:"size:512"`
	OfflineContent string `json:"offline_content" gorm:"size:512"`
	Name string `json:"name" gorm:"size:64"`
	LastOnline time.Time `json:"last_online"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

func (setting *AdminChatSetting) GetOfflineMsg(uid int64, sessionId uint64) *Message  {
	offlineMsg := &Message{
		UserId:     uid,
		AdminId:    setting.AdminId,
		Type:       TypeText,
		Content:    setting.OfflineContent,
		ReceivedAT: time.Now().Unix(),
		Source:     SourceAdmin,
		SessionId:  sessionId,
		ReqId:      databases.GetSystemReqId(),
		IsRead:     false,
	}
	return offlineMsg
}
