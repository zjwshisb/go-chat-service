package models


type Message struct {
	Id uint64 `gorm:"primaryKey"`
	UserId int64 `gorm:"index"`
	ServiceId int64 `gorm:"index"`
	Type string `gorm:"size:16"`
	Content string `gorm:"size:1024"`
	CreatedAT int
	ReqId int64
	IsServer bool
}

