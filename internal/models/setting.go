package models

import (
	"time"
	"ws/internal/json"
)

type Setting struct {
	Id          string    `gorm:"id;primaryKey" json:"id"`
	Value       string    `gorm:"value" json:"value"`
	Name        string    `gorm:"name" json:"name"`
	Description string    `gorm:"description" json:"description"`
	Options []json.Options `gorm:"options" json:"options"`
	CreatedAt   time.Time `gorm:"created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"updated_at" json:"updated_at"`
}


