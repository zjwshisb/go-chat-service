package models

import (
	"gorm.io/gorm/clause"
	"ws/app/databases"
	"ws/app/resource"
)

const (
	TypeImage    = "image"
	TypeText     = "text"
	TypeNavigate = "navigator"
	TypeNotice   = "notice"
	SourceUser   = 0
	SourceAdmin  = 1
	SourceSystem = 2
)

type Message struct {
	Id         uint64 `gorm:"primaryKey"`
	UserId     int64  `gorm:"index" mapstructure:"user_id"`
	AdminId    int64  `gorm:"index"`
	Type       string `gorm:"size:16" mapstructure:"type"`
	Content    string `gorm:"size:1024" mapstructure:"content"`
	ReceivedAT int64
	GroupId    int64            `gorm:"group_id"`
	SendAt     int64  `gorm:"send_at"`
	Source     int8   `gorm:"source"`
	SessionId  uint64 `gorm:"session_id"`
	ReqId      string  `gorm:"index" mapstructure:"req_id"`
	IsRead     bool   `gorm:"bool"`
	Admin      *Admin `gorm:"foreignKey:admin_id"`
	User       *User  `gorm:"foreignKey:user_id"`
}

func (message *Message) Save() {
	databases.Db.Omit(clause.Associations).Save(message)
}

func (message *Message) GetUser() *User  {
	if message.User == nil {
		user := &User{}
		_ = databases.Db.Model(message).Association("User").Find(user)
		message.User = user
	}
	return message.User
}
func (message *Message) GetAdmin() *Admin  {
	if message.Admin == nil {
		admin := &Admin{}
		_ = databases.Db.Model(message).Association("Admin").Find(admin)
		message.Admin = admin
	}
	return message.Admin
}
func (message *Message) GetAdminName() string {
	switch message.Source {
	case SourceAdmin:
		return message.GetAdmin().GetChatName()
	case SourceSystem:
		setting := &ChatSetting{}
		databases.Db.Where("name = ?", SystemName).
			Where("group_id = ?", message.GroupId).First(setting)
		return setting.Value
	}
	return ""
}
func (message *Message) GetAvatar() (avatar string) {
	switch message.Source {
	case SourceUser:
		avatar = message.GetUser().GetAvatarUrl()
	case SourceAdmin:
		avatar = message.GetAdmin().GetAvatarUrl()
	case SourceSystem:
		setting := &ChatSetting{}
		databases.Db.Where("name = ?", SystemAvatar).
			Where("group_id = ?", message.GroupId).First(setting)
		return setting.Value
	}
	return
}
func (message *Message) ToJson() *resource.Message {
	return &resource.Message{
		Id:         message.Id,
		UserId:     message.UserId,
		AdminId:    message.AdminId,
		AdminName:  message.GetAdminName(),
		Type:       message.Type,
		Content:    message.Content,
		ReceivedAT: message.ReceivedAT,
		Source:     message.Source,
		ReqId:      message.ReqId,
		IsSuccess:  true,
		IsRead:     message.IsRead,
		Avatar:     message.GetAvatar(),
	}
}

