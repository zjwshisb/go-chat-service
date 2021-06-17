package models

import "time"

type SubScribeMessage struct {
	ID int64 `gorm:"primaryKey,autoIncrement"`
	TemplateId string `gorm:"template_id"`
	UserId int64 `gorm:"user_id"`
	IsUseD bool `gorm:"is_used"`
	CreatedAt time.Time
	UpdatedAt time.Time
}