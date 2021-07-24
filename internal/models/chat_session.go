package models

type ChatSession struct {
	Id uint64 `gorm:"primaryKey" json:"id"`
	UserId int64 `gorm:"index"`
	QueriedAt int64 `json:"queried_at"`
	AcceptedAt int64 `json:"accepted_at"`
	BrokeAt int64	`json:"broke_at"`
	ServiceId int64 `gorm:"index"`
	BackendUser *BackendUser `gorm:"foreignKey:service_id"`
	User       *User        `gorm:"foreignKey:user_id"`
}
