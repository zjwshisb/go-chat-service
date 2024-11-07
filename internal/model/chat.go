package model

import "github.com/gogf/gf/v2/os/gtime"

type ChatOnlineCount struct {
	Admin   uint
	User    uint
	Waiting uint
}

type ChatAction struct {
	Data   any    `json:"data"`
	Time   int64  `json:"time"`
	Action string `json:"action"`
}

type ChatWaitingUser struct {
	Username     string              `json:"username"`
	Avatar       string              `json:"avatar"`
	UserId       uint                `json:"id"`
	LastTime     *gtime.Time         `json:"last_time"`
	Messages     []ChatSimpleMessage `json:"messages"`
	MessageCount uint                `json:"message_count"`
	Description  string              `json:"description"`
	SessionId    uint                `json:"session_id"`
}

type ChatSimpleMessage struct {
	Type    string      `json:"type"`
	Time    *gtime.Time `json:"time"`
	Content string      `json:"content"`
}

type ChatCustomerAdmin struct {
	Avatar        string `json:"avatar"`
	Username      string `json:"username"`
	Online        bool   `json:"online"`
	Id            uint   `json:"id"`
	AcceptedCount uint   `json:"accepted_count"`
	Platform      string `json:"platform"`
}

type ChatTransfer struct {
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

type ChatMessage struct {
	Id         uint        `json:"id"`
	UserId     uint        `json:"user_id"`
	AdminId    uint        `json:"admin_id"`
	AdminName  string      `json:"admin_name"`
	Type       string      `json:"type"`
	Content    string      `json:"content"`
	ReceivedAT *gtime.Time `json:"received_at"`
	Source     uint        `json:"source"`
	ReqId      string      `json:"req_id"`
	IsSuccess  bool        `json:"is_success"`
	IsRead     bool        `json:"is_read"`
	Avatar     string      `json:"avatar"`
	Username   string      `json:"username"`
}

type ChatUser struct {
	Id           uint         `json:"id"`
	Username     string       `json:"username"`
	LastChatTime *gtime.Time  `json:"last_chat_time"`
	Disabled     bool         `json:"disabled"`
	Online       bool         `json:"online"`
	LastMessage  *ChatMessage `json:"last_message"`
	Unread       uint         `json:"unread"`
	Avatar       string       `json:"avatar"`
	Platform     string       `json:"platform"`
}

type ChatSimpleUser struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
}

type ChatSession struct {
	Id          uint        `json:"id"`
	UserId      uint        `json:"-"`
	QueriedAt   *gtime.Time `json:"queried_at"`
	AcceptedAt  *gtime.Time `json:"accepted_at"`
	BrokenAt    *gtime.Time `json:"broken_at"`
	CanceledAt  *gtime.Time `json:"canceled_at"`
	AdminId     uint        `json:"admin_id"`
	UserName    string      `json:"username"`
	AdminName   string      `json:"admin_name"`
	TypeLabel   string      `json:"type_label"`
	Status      string      `json:"status"`
	StatusLabel string      `json:"status_label"`
	Rate        uint        `json:"rate"`
}
