package chat

import "time"

type WaitingUser struct {
	Username     string          `json:"username"`
	Avatar       string          `json:"avatar"`
	UserId       uint            `json:"id"`
	LastTime     int64           `json:"last_time"`
	Messages     []SimpleMessage `json:"messages"`
	MessageCount uint            `json:"message_count"`
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
	AcceptedCount uint   `json:"accepted_count"`
	Platform      string `json:"platform"`
}

type Transfer struct {
	Id            uint      `json:"id"`
	FromSessionId uint      `json:"from_session_id"`
	ToSessionId   uint      `json:"to_session_id"`
	UserId        uint      `json:"user_id"`
	Remark        string    `json:"remark"`
	FromAdminName string    `json:"from_admin_name"`
	ToAdminName   string    `json:"to_admin_name"`
	Username      string    `json:"username"`
	CreatedAt     time.Time `json:"created_at"`
	AcceptedAt    time.Time `json:"accepted_at"`
	CanceledAt    time.Time `json:"canceled_at"`
}
