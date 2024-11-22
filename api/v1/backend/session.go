package backend

import (
	"gf-chat/api"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/gogf/gf/v2/frame/g"
)

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

type SessionListReq struct {
	g.Meta `path:"/chat-sessions" tags:"客服对话" method:"get" summary:"客户对话列表"`
	api.Paginate
	QueriedAt map[string]string `p:"queried" dc:"消息类型"`
	Username  string            `dc:"用户手机号"`
	AdminName string            `dc:"客服名称"`
	Status    string            `dc:"状态"`
}

type SessionListRes struct {
	Items []ChatSession `json:"items"`
	Total int           `json:"total"`
}

type SessionCancelReq struct {
	g.Meta `path:"/chat-sessions/:id/cancel" tags:"客服对话" method:"post" summary:"取消客服对话"`
}

type SessionCloseReq struct {
	g.Meta `path:"/chat-sessions/:id/close" tags:"客服对话" method:"post" summary:"关闭客服对话"`
}

type SessionDetailReq struct {
	g.Meta `path:"/chat-sessions/:id" tags:"客服对话" method:"get" summary:"获取客服对话详情"`
}

type SessionDetailRes struct {
	Messages []ChatMessage `json:"messages"`
	Session  ChatSession   `json:"session"`
	Total    int           `json:"total"`
}
