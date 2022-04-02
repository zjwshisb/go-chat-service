package mq

import (
	"encoding/json"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"strings"
	"ws/app/log"
)

const TypeMessage = "message"
const TypeWaitingUser = "waiting-user"
const TypeAdmin = "admin"
const TypeOtherLogin = "other-login"
const TypeMoreThanOne = "more-than-one"
const TypeTransfer = "admin-transfer"
const TypeWaitingUserCount = "waiting-user-count"
const TypeUpdateSetting = "update-admin-setting"
const TypeUserOnline = "user-online"
const TypeUserOffline = "user-offline"

type MessageQueue interface {
	// Publish 消息
	Publish(channel string, p *Payload) error
	// Subscribe 消息
	Subscribe(channel string) SubScribeChannel
}

type SubScribeChannel interface {
	// ReceiveMessage 接收消息
	ReceiveMessage() gjson.Result
	Close()
}

type Payload struct {
	Types string      `json:"types"`
	Data  interface{} `json:"data"`
}

func (payload *Payload) MarshalBinary() ([]byte, error) {
	return json.Marshal(payload)
}

func NewMessagePayload(mid uint64) *Payload {
	return &Payload{
		Types: TypeMessage,
		Data:  mid,
	}
}

var mq MessageQueue

func Setup() {
	switch strings.ToLower(viper.GetString("App.Mq")) {
	case "rabbitmq":
		mq = newRabbitMq()
	case "redis":
		mq = newRedisMq()
	default:
		mq = newRedisMq()
	}
}

func Publish(channel string, p *Payload) error {
	log.Log.WithField("a-type", "publish/subscribe").
		WithField("b-type", "publish").
		Infof("<channel:%s><types:%s><data:%v>", channel, p.Types, p.Data)
	return mq.Publish(channel, p)
}

func Subscribe(channel string) SubScribeChannel {
	return mq.Subscribe(channel)
}
