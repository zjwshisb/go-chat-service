package chat

type Action struct {
	Data   any    `json:"data"`
	Time   int64  `json:"time"`
	Action string `json:"action"`
}
