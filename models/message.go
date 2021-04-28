package models

type Message struct {
	Id uint64 `gorm:"primaryKey"`
	UserId int64 `gorm:"index" mapstructure:"user_id"`
	ServiceId int64 `gorm:"index"`
	Type string `gorm:"size:16" mapstructure:"type"`
	Content string `gorm:"size:1024" mapstructure:"content"`
	ReceivedAT int64 `json:"received_at"`
	SendAt int64
	IsServer bool `gorm:"is_server"`
	ReqId int64 `gorm:"index" mapstructure:"req_id"`
	IsRead bool `gorm:"bool" json:"is_read"`
	ServerUser ServerUser `gorm:"foreignKey:service_id"`
	User User `gorm:"foreignKey:user_id"`
}
