package modules

import "encoding/json"

type Action struct {
	ReqId int64 `json:"req_id"`
	Data interface{} `json:"data"`
	Time int64	`json:"time"`
	Action string `json:"action"`
}

func (action *Action) Marshal() (b []byte) {
	b, err := json.Marshal(action)
	if err != nil {
		b = nil
	}
	return
}

func (action *Action) UnMarshal(b []byte) (err error) {
	err =  json.Unmarshal(b, action)
	return
}