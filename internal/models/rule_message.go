package models

import "time"

const (
	TypeImage = "image"
	TypeText = "text"
	TypeNavigate = "navigator"

	MatchTypeAll = "all"
	MatchTypePart = "part"
	ReplyTypeMessage = "message"
	ReplyTypeTransfer = "transfer"
)

type AutoMessage struct {
	ID uint `gorm:"column:id;primaryKey" json:"id"`
	Name string `gorm:"size:255" json:"name"`
	Type string  `gorm:"size:255" json:"type"`
	Content string `gorm:"type:text" json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


type AutoRule struct {
	ID uint
	Name string `gorm:"size:255"`
	Match string `gorm:"size:32"`
	MatchType string `gorm:"size:20"`
	ReplyType string `gorm:"size:20"`
	Sort uint
	Count uint `gorm:"not null;default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}