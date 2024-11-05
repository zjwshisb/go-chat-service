package chat

import "time"

type WaitingUser struct {
	Username     string          `json:"username"`
	Avatar       string          `json:"avatar"`
	UserId       int             `json:"id"`
	LastTime     int64           `json:"last_time"`
	Messages     []SimpleMessage `json:"messages"`
	MessageCount int             `json:"message_count"`
	Description  string          `json:"description"`
	SessionId    uint64          `json:"session_id"`
}

type SimpleMessage struct {
	Type    string `json:"type"`
	Time    int64  `json:"time"`
	Content string `json:"content"`
}

type CustomerAdmin struct {
	Avatar        string `json:"avatar"`
	Username      string `json:"username"`
	Online        bool   `json:"online"`
	Id            uint   `json:"id"`
	AcceptedCount int    `json:"accepted_count"`
	Platform      string `json:"platform"`
}

type Transfer struct {
	Id            int       `json:"id"`
	FromSessionId uint64    `json:"from_session_id"`
	ToSessionId   uint64    `json:"to_session_id"`
	UserId        int       `json:"user_id"`
	Remark        string    `json:"remark"`
	FromAdminName string    `json:"from_admin_name"`
	ToAdminName   string    `json:"to_admin_name"`
	Username      string    `json:"username"`
	CreatedAt     time.Time `json:"created_at"`
	AcceptedAt    int64     `json:"accepted_at"`
	CanceledAt    int64     `json:"canceled_at"`
}
