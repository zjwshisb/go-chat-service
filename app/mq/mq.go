package mq

import (
	"encoding/json"
	"strings"
	"ws/configs"
)

const TypeMessage = "message"
const TypeWaitingUser = "waiting-user"
const TypeAdmin = "admin"
const TypeAdminLogin = "admin-login"

type MessageQueue interface {
	// 发布
	Publish(channel string, p *Payload) error
	// 订阅
	Subscribe(channel string) SubScribeChannel
}

type SubScribeChannel interface {
	// 接收消息
	ReceiveMessage() (*Payload, error)
	Close()
}


type Payload struct {
	Types string `json:"types"`
	Data interface{} `json:"data"`
}

func (payload *Payload) MarshalBinary() ([]byte, error) {
	return json.Marshal(payload)
}

func NewMessagePayload(mid uint64) *Payload  {
	return &Payload{
		Types: TypeMessage,
		Data:  mid,
	}
}

var mq MessageQueue

func init()  {
	switch strings.ToLower(configs.App.Mq) {
	case "rabbitmq":
		mq = newRabbitMq()
	case "redis":
		mq = newRedisMq()
	default:
		mq = newRedisMq()
	}
}

func Mq() MessageQueue {
	return mq
}