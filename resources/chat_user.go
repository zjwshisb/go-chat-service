package resources

import "ws/models"

type ChatUser struct {
	ID        int64      `json:"id"`
	Username  string     `json:"username"`
	LastChatTime int64  `json:"last_chat_time"`
	Disabled bool `json:"disabled"`
	Online bool `json:"online"`
	Messages []*Message `json:"messages"`
	Unread int `json:"unread"`
}

func NewChatUser(user models.User) *ChatUser {
	return &ChatUser{
		ID: user.ID,
		Username: user.Username,
		LastChatTime: 0,
		Messages: make([]*Message, 0),
	}
}