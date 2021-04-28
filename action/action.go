package action

import (
	"encoding/json"
	"errors"
	"github.com/mitchellh/mapstructure"
	"time"
	"ws/models"
	"ws/resources"
)

const (
	ReceiptAction = "receipt"
	PingAction = "ping"
	UserOnLineAction = "user-online"
	UserOffLineAction = "user-offline"
	UserWaitingCountAction = "waiting-users"
	SendMessageAction = "send-message"
	ReceiveMessageAction = "receive-message"
)

type Action struct {
	Data interface{} `json:"data"`
	Time int64	`json:"time"`
	Action string `json:"action"`
}

func (action *Action) Marshal() (b []byte, err error) {
	if action.Action == PingAction {
		return []byte(""), nil
	}
	if action.Action == ReceiveMessageAction {
		msg, ok := action.Data.(*models.Message)
		if !ok {
			err = errors.New("param error")
			return
		}
		b, err = json.Marshal(Action{
			Time: action.Time,
			Action: action.Action,
			Data: resources.NewMessage(*msg),
		})
		return
	}
	b, err = json.Marshal(action)
	return
}

func (action *Action) UnMarshal(b []byte) (err error) {
	err =  json.Unmarshal(b, action)
	return
}

func (action *Action) GetMessage() (message *models.Message,err error)  {
	if action.Action == SendMessageAction {
		message = &models.Message{}
		err = mapstructure.Decode(action.Data, message)
	} else {
		err = errors.New("无效的action")
	}
	return

}
func NewReceiveAction (msg *models.Message) *Action {
	return &Action{
		Action: ReceiveMessageAction,
		Time: time.Now().Unix(),
		Data: msg,
	}
}
func NewReceipt(msg *models.Message) (act *Action) {
	data := make(map[string]interface{})
	data["user_id"] = msg.UserId
	data["req_id"] = msg.ReqId
	act = &Action{
		Action: ReceiptAction,
		Time: time.Now().Unix(),
		Data: data,
	}
	return
}
func NewUserOnline(uid int64) *Action {
	data := make(map[string]interface{})
	data["user_id"] = uid
	return &Action{
		Action: UserOnLineAction,
		Time: time.Now().Unix(),
		Data: data,
	}
}
func NewUserOffline(uid int64) *Action {
	data := make(map[string]interface{})
	data["user_id"] = uid
	return &Action{
		Action: UserOffLineAction,
		Time: time.Now().Unix(),
		Data: data,
	}
}
func NewPing() *Action {
	return &Action{
		Action: PingAction,
		Time: time.Now().Unix(),
	}
}
func NewWaitingUsers(i interface{}) *Action {
	return &Action{
		Action: UserWaitingCountAction,
		Time: time.Now().Unix(),
		Data: i,
	}
}


