package models

import (
	"gorm.io/gorm/clause"
	"ws/app/databases"
	"ws/app/resource"
	"ws/app/util"
	"ws/configs"
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
	SendAt     int64  `gorm:"send_at"`
	Source     int8   `gorm:"source"`
	SessionId  uint64 `gorm:"session_id"`
	ReqId      int64  `gorm:"index" mapstructure:"req_id"`
	IsRead     bool   `gorm:"bool"`
	Admin      *Admin `gorm:"foreignKey:admin_id"`
	User       *User  `gorm:"foreignKey:user_id"`
}

func (message *Message) Save() {
	databases.Db.Omit(clause.Associations).Save(message)
}

func (message *Message) GetAdminName() string {
	switch message.Source {
	case SourceAdmin:
		if message.Admin != nil {
			return message.Admin.GetChatName()
		}
	case SourceSystem:
		return configs.App.SystemChatName
	}
	return ""
}
func (message *Message) GetAvatar() (avatar string) {
	switch message.Source {
	case SourceUser:
		if message.User != nil {
			avatar = message.User.GetAvatarUrl()
		}
	case SourceAdmin:
		if message.Admin != nil {
			avatar = message.Admin.GetAvatarUrl()
		}
	case SourceSystem:
		avatar = util.SystemAvatar()
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

