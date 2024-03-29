package models

import (
	"encoding/json"
	"time"
	"ws/app/resource"
)


const (
	IsAutoTransfer = "is-auto-transfer"
	MinuteToBreak = "minute-to-break"
	SystemName = "system-name"
	SystemAvatar = "system-avatar"
)

type ChatSetting struct {
	Id int64
	Name string `gorm:"size:255"`
	Title string `gorm:"size:255"`
	GroupId int64 `gorm:"index"`
	Value string  `gorm:"size:255"`
	Options string `gorm:"size:1024"`
	Type string `gorm:"size:16"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func (setting *ChatSetting) ToJson() *resource.ChatSetting {
	var o = make([]map[string]string, 0)
	_ = json.Unmarshal([]byte(setting.Options), &o)
	return &resource.ChatSetting{
		Id:      setting.Id,
		Name:    setting.Name,
		Title:   setting.Title,
		Value:   setting.Value,
		Options: o,
		Type: setting.Type,
	}
}