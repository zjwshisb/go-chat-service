package resources

type WaitingUser struct {
	Username string `json:"username"`
	Avatar string `json:"avatar"`
	Id int64 `json:"id"`
	LastMessage string `json:"last_message"`
	LastTime int64 `json:"last_time"`
	MessageCount int `json:"message_count"`
	Description string `json:"description"`
}