package models

type QueryRecord struct {
	Id uint64 `gorm:"primaryKey"`
	UserId int64 `gorm:"index"`
	QueriedAt int64
	AcceptedAt int64
	BrokeAt int64
	ServiceId int64 `gorm:"index"`
}
