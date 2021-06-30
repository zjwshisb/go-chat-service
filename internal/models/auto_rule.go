package models

import (
	"time"
	"ws/internal/chat"
	"ws/internal/util"
)

const (
	MatchTypeAll  = "all"
	MatchTypePart = "part"

	MatchEnter             = "enter"
	MatchServiceAllOffLine = "u-offline"

	ReplyTypeMessage  = "message"
	ReplyTypeTransfer = "transfer"
)

type AutoRule struct {
	ID        uint         `json:"id"`
	Name      string       `gorm:"size:255" json:"name"`
	Match     string       `gorm:"size:32" json:"match"`
	MatchType string       `gorm:"size:20" json:"match_type"`
	ReplyType string       `gorm:"size:20" json:"reply_type"`
	MessageId uint         `gorm:"index" json:"message_id"`
	IsSystem  uint8        `gorm:"is_system" json:"-"`
	Sort      uint8        `gorm:"sort" json:"sort"`
	IsOpen    bool         `gorm:"is_open" json:"is_open"`
	Count     uint         `gorm:"not null;default:0" json:"count"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	Message   *AutoMessage `json:"message" gorm:"foreignKey:MessageId"`
}

func (rule *AutoRule) GetMessages(uid int64) (message *Message) {
	if rule.Message != nil {
		message = &Message{
			UserId:     uid,
			ServiceId:  0,
			Type:       rule.Message.Type,
			Content:    rule.Message.Content,
			ReceivedAT: time.Now().Unix(),
			SendAt:     0,
			Source:     SourceSystem,
			ReqId:      util.CreateReqId(),
			IsRead:     true,
			Avatar:     chat.SystemAvatar(),
		}
	}
	return
}
