package chat

import (
	"github.com/gogf/gf/v2/os/gtime"
)

type WaitingUser struct {
	Username     string          `json:"username"`
	Avatar       string          `json:"avatar"`
	UserId       uint            `json:"id"`
	LastTime     *gtime.Time     `json:"last_time"`
	Messages     []SimpleMessage `json:"messages"`
	MessageCount uint            `json:"message_count"`
	Description  string          `json:"description"`
	SessionId    uint            `json:"session_id"`
}

type SimpleMessage struct {
	Type    string      `json:"type"`
	Time    *gtime.Time `json:"time"`
	Content string      `json:"content"`
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
	Id            uint        `json:"id"`
	FromSessionId uint        `json:"from_session_id"`
	ToSessionId   uint        `json:"to_session_id"`
	UserId        uint        `json:"user_id"`
	Remark        string      `json:"remark"`
	FromAdminName string      `json:"from_admin_name"`
	ToAdminName   string      `json:"to_admin_name"`
	Username      string      `json:"username"`
	CreatedAt     *gtime.Time `json:"created_at"`
	AcceptedAt    *gtime.Time `json:"accepted_at"`
	CanceledAt    *gtime.Time `json:"canceled_at"`
}
