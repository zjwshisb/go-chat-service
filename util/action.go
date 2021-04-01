package util

import (
	"encoding/json"
	"time"
)

type Action struct {
	ReqId int64 `json:"req_id"`
	Data map[string]interface{} `json:"data"`
	Time int64	`json:"time"`
	Action string `json:"action"`
}

func (action *Action) Marshal() (b []byte, err error) {
	b, err = json.Marshal(action)
	return
}

func (action *Action) UnMarshal(b []byte) (err error) {
	err =  json.Unmarshal(b, action)
	return
}
func NewReceiptAction(action  Action) *Action  {
	data := make(map[string]interface{})
	userId, ok := action.Data["user_id"]
	if ok {
		data["user_id"] = userId
	}
	return &Action{
		ReqId: action.ReqId,
		Action: "receipt",
		Time: time.Now().Unix(),
		Data: data,
	}
}