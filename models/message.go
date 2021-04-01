package models

import (
	"github.com/mitchellh/mapstructure"
	"time"
	"ws/util"
)

type Message struct {
	Id uint64 `gorm:"primaryKey"`
	UserId int64 `gorm:"index" mapstructure:"user_id"`
	ServiceId int64 `gorm:"index"`
	Type string `gorm:"size:16" mapstructure:"type"`
	Content string `gorm:"size:1024" mapstructure:"content"`
	CreatedAT int64
	IsServer bool `gorm:"is_server"`
}

func NewFromAction(action util.Action) (message *Message,err error) {
	message = &Message{}
	err = mapstructure.Decode(action.Data, message)
	if err == nil {
		message.CreatedAT = time.Now().Unix()
	}
	return
}
