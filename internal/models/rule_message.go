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
	ID uint `gorm:"column:id;primaryKey"`
	Name string `gorm:"size:255"`
	Type string  `gorm:"size:255"`
	Content string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
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