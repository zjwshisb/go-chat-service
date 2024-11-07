package backend

import (
	"gf-chat/internal/model/chat"

	"github.com/gogf/gf/v2/frame/g"
)

type ListReq struct {
	g.Meta    `path:"/chat-sessions" tags:"客服对话" method:"get" summary:"客户对话列表"`
	PageSize  int               `d:"20" json:"pageSize" v:"max:100"`
	Current   int               `d:"1" dc:"页码" json:"current"`
	QueriedAt map[string]string `p:"queried" dc:"消息类型"`
	Username  string            `dc:"用户手机号"`
	AdminName string            `dc:"客服名称"`
	Status    string            `dc:"状态"`
}

type ListRes struct {
	Items []chat.Session `json:"items"`
	Total int            `json:"total"`
}

type CancelReq struct {
	g.Meta `path:"/chat-sessions/:id/cancel" tags:"客服对话" method:"post" summary:"取消客服对话"`
}

type CloseReq struct {
	g.Meta `path:"/chat-sessions/:id/close" tags:"客服对话" method:"post" summary:"关闭客服对话"`
}

type DetailReq struct {
	g.Meta `path:"/chat-sessions/:id" tags:"客服对话" method:"get" summary:"获取客服对话详情"`
}

type DetailRes struct {
	Messages []chat.Message `json:"messages"`
	Session  chat.Session   `json:"session"`
	Total    int            `json:"total"`
}
