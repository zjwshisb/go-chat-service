package models

import "time"

type AutoMessageJson struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	RulesCount int       `json:"rules_count"`
}

type AutoMessage struct {
	ID        uint   `gorm:"column:id;primaryKey"`
	Name      string `gorm:"size:255"`
	Type      string `gorm:"size:255"`
	Content   string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Rules     []*AutoRule `gorm:"foreignKey:message_id"`
}

func (message *AutoMessage) ToJson() *AutoMessageJson {
	count := 0
	if message.Rules != nil {
		count = len(message.Rules)
	}
	return &AutoMessageJson{
		ID:         message.ID,
		Name:       message.Name,
		Type:       message.Type,
		Content:    message.Content,
		CreatedAt:  message.CreatedAt,
		UpdatedAt:  message.UpdatedAt,
		RulesCount: count,
	}
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
