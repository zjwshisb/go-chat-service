package chat

type User struct {
	Id           int      `json:"id"`
	Username     string   `json:"username"`
	LastChatTime int64    `json:"last_chat_time"`
	Disabled     bool     `json:"disabled"`
	Online       bool     `json:"online"`
	LastMessage  *Message `json:"last_message"`
	Unread       int      `json:"unread"`
	Avatar       string   `json:"avatar"`
	Platform     string   `json:"platform"`
}

type SimpleUser struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}
