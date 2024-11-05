package chat

import "github.com/gogf/gf/v2/os/gtime"

type Message struct {
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
