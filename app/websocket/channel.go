package websocket

import "encoding/json"

const TypeMessage = "message"
const TypeWaitingUser = "waiting-user"
const TypeAdmin = "admin"

type payload struct {
	Types string `json:"types"`
	Data interface{} `json:"data"`
}

func (payload *payload) MarshalBinary() ([]byte, error) {
	return json.Marshal(payload)
}

func newMessagePayload(mid uint64) *payload  {
	return &payload{
		Types: TypeMessage,
		Data:  mid,
	}
}

