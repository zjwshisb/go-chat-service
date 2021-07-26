package json

import "ws/internal/models"

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	LastChatTime int64  `json:"last_chat_time"`
	Disabled bool       `json:"disabled"`
	Online bool         `json:"online"`
	Messages []*models.MessageJson `json:"messages"`
	Unread int          `json:"unread"`
}
