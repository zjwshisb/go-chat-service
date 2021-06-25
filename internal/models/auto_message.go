package models

import "time"

const (
	TypeImage = "image"
	TypeText = "text"
	TypeNavigate = "navigator"
)

type AutoMessage struct {
	ID uint `gorm:"column:id;primaryKey" json:"id"`
	Name string `gorm:"size:255" json:"name"`
	Type string  `gorm:"size:255" json:"type"`
	Content string `gorm:"type:text" json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


