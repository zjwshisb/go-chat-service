package util

import "encoding/json"

type Action struct {
	ReqId int64 `json:"req_id"`
	Data interface{} `json:"data"`
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