package models

import "time"

const (
	MatchTypeAll = "all"
	MatchTypePart = "part"

	ReplyTypeMessage = "message"
	ReplyTypeTransfer = "transfer"
)

type AutoRule struct {
	ID uint `json:"id"`
	Name string `gorm:"size:255" json:"name"`
	Match string `gorm:"size:32" json:"match"`
	MatchType string `gorm:"size:20" json:"match_type"`
	ReplyType string `gorm:"size:20" json:"reply_type"`
	MessageId uint `gorm:"index" json:"message_id"`
	Sort uint `json:"sort"`
	Count uint `gorm:"not null;default:0" json:"count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
