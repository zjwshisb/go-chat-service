package models

import (
	"encoding/json"
	"time"
)

const (
	messageAction = "message"
	receiptAction = "receipt"
	userOnLineAction = "user-online"
	userOffLineAction = "user-offline"
	userWaitingCountAction = "user-waiting-count"
)

type Action struct {
	ReqId int64 `json:"req_id"`
	Data map[string]interface{} `json:"data"`
	Time int64	`json:"time"`
	Action string `json:"action"`
	Message *Message `json:"-"`
}

func (action *Action) Marshal() (b []byte, err error) {
	b, err = json.Marshal(action)
	return
}

func (action *Action) UnMarshal(b []byte) (err error) {
	err =  json.Unmarshal(b, action)
	return
}
func NewReceiptAction(action *Action) *Action {
	data := make(map[string]interface{})
	userId, ok := action.Data["user_id"]
	if ok {
		data["user_id"] = userId
	}
	return &Action{
		ReqId: action.ReqId,
		Action: receiptAction,
		Time: time.Now().Unix(),
		Data: data,
	}
}
func NewUserOnlineAction(uid int64) *Action {
	data := make(map[string]interface{})
	data["user_id"] = uid
	return &Action{
		ReqId: 11,
		Action: userOnLineAction,
		Time: time.Now().Unix(),
		Data: data,
	}
}
func NewUserOfflineAction(uid int64) *Action {
	data := make(map[string]interface{})
	data["user_id"] = uid
	return &Action{
		ReqId: 11,
		Action: userOffLineAction,
		Time: time.Now().Unix(),
		Data: data,
	}
}
func NewPingAction()  *Action {
	return &Action{
		ReqId: 12,
		Action: "ping",
		Time: time.Now().Unix(),
	}
}
func NewUserWaitingCountAction(count int) *Action {
	data := make(map[string]interface{})
	data["count"] = count
	return &Action{
		ReqId: 12,
		Action: userWaitingCountAction,
		Time: time.Now().Unix(),
		Data: data,
	}
}
func NewServiceOnlineListAction(data map[string]interface{}) *Action {
	return &Action{
		Action: "service_online_count",
		Time: time.Now().Unix(),
		Data: data,
	}
}