package models

import "ws/app/util"

const (
	TypeImage = "image"
	TypeText = "text"
	TypeNavigate = "navigator"

	SourceUser = 0
	SourceAdmin = 1
	SourceSystem = 2
)

type Message struct {
	Id         uint64 `gorm:"primaryKey"`
	UserId     int64 `gorm:"index" mapstructure:"user_id"`
	AdminId  int64 `gorm:"index"`
	Type       string `gorm:"size:16" mapstructure:"type"`
	Content    string `gorm:"size:1024" mapstructure:"content"`
	ReceivedAT int64 `json:"received_at"`
	SendAt     int64 `json:"send_at" gorm:"send_at"`
	Source   int8        `gorm:"source"`
	SessionId uint64 `gorm:"session_id"`
	ReqId      int64       `gorm:"index" mapstructure:"req_id"`
	IsRead     bool        `gorm:"bool" json:"is_read"`
	Admin *Admin `gorm:"foreignKey:admin_id"`
	User       *User        `gorm:"foreignKey:user_id"`
	Avatar    string       `gorm:"-"`
}

func (message *Message) GetAvatar() (avatar string)  {
	if message.Avatar != "" {
		return message.Avatar
	}
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

func (message *Message) ToJson() *MessageJson {
	return &MessageJson{
		Id:         message.Id,
		UserId:     message.UserId,
		AdminId:  message.AdminId,
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

type MessageJson struct {
	Id         uint64 `json:"id"`
	UserId     int64  `json:"user_id"`
	AdminId  int64  `json:"admin_id"`
	Type       string `json:"type"`
	Content    string `json:"content"`
	ReceivedAT int64  `json:"received_at"`
	Source   int8   `json:"source"`
	ReqId      int64  `json:"req_id"`
	IsSuccess  bool   `json:"is_success"`
	IsRead     bool   `json:"is_read"`
	Avatar     string `json:"avatar"`
}
