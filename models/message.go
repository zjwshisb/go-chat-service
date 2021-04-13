package models

type Message struct {
	Id uint64 `gorm:"primaryKey" json:"id"`
	UserId int64 `gorm:"index" mapstructure:"user_id" json:"user_id"`
	ServiceId int64 `gorm:"index" json:"service_id"`
	Type string `gorm:"size:16" mapstructure:"type" json:"type"`
	Content string `gorm:"size:1024" mapstructure:"content" json:"content"`
	ReceivedAT int64 `json:"received_at"`
	SendAt int64 `json:"-"`
	IsServer bool `gorm:"is_server" json:"is_server"`
}

