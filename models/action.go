package models

import (
	"encoding/json"
	"time"
	"ws/util"
)

const (
	messageAction = "message"
	receiptAction = "receipt"
	pingAction = "ping"
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
	if action.Action == pingAction {
	} else {
		b, err = json.Marshal(action)
	}
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
		ReqId: util.CreateReqId(),
		Action: userOnLineAction,
		Time: time.Now().Unix(),
		Data: data,
	}
}
func NewUserOfflineAction(uid int64) *Action {
	data := make(map[string]interface{})
	data["user_id"] = uid
	return &Action{
		ReqId: util.CreateReqId(),
		Action: userOffLineAction,
		Time: time.Now().Unix(),
		Data: data,
	}
}
func NewPingAction()  *Action {
	return &Action{
		ReqId: util.CreateReqId(),
		Action: pingAction,
		Time: time.Now().Unix(),
	}
}
func NewUserWaitingCountAction(count int) *Action {
	data := make(map[string]interface{})
	data["count"] = count
	return &Action{
		ReqId: util.CreateReqId(),
		Action: userWaitingCountAction,
		Time: time.Now().Unix(),
		Data: data,
	}
}
func NewServiceOnlineListAction(data map[string]interface{}) *Action {
	return &Action{
		ReqId: util.CreateReqId(),
		Action: "service_online_list",
		Time: time.Now().Unix(),
		Data: data,
	}
}