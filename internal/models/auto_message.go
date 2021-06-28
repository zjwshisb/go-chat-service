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

func (message *AutoMessage) TypeLabel() string {
	switch message.Type {
	case TypeText:
		return "文本"
	case TypeNavigate:
		return "导航卡片"
	case TypeImage:
		return "图片"
	default:
		return "未知类型"
	}
}


