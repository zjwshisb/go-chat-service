package websocket

import "encoding/json"

type payload struct {
	Types string `json:"types"`
	Data interface{} `json:"data"`
}

func (payload *payload) MarshalBinary() ([]byte, error) {
	return json.Marshal(payload)
}

func (payload *payload) UnmarshalBinary(data []byte) error  {
	return json.Unmarshal(data, payload)
}
