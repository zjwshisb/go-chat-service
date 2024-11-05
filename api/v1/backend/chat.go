package backend

import (
	"gf-chat/internal/model/chat"
	"github.com/gogf/gf/v2/frame/g"
)

type ChatUserInfoReq struct {
	g.Meta `path:"/ws/chat-user/:id" tags:"后台客服面板" method:"get" summary:"获取用户信息"`
}
type ChatSmsNoticeReq struct {
	g.Meta `path:"/ws/sms-notice" tags:"后台客服面板" method:"post" summary:"发送短信提醒"`
	Uid    int `json:"uid" v:"required"`
}
type ChatAcceptReq struct {
	g.Meta `path:"/ws/chat-user" tags:"后台客服面板" method:"post" summary:"接入用户"`
	Sid    uint64 `p:"sid" json:"sid"`
}
type ChatRemoveReq struct {
	g.Meta `path:"/ws/chat-user/:id" tags:"后台客服面板" method:"delete" summary:"移除用户"`
}
type ChatRemoveAllReq struct {
	g.Meta `path:"/ws/chat-user" tags:"后台客服面板" method:"delete" summary:"移除所有失效用户"`
}
type ChatRemoveAllRes struct {
	Ids []int `json:"ids" d:"移除的用户Id"`
}

type ChatReadReq struct {
	g.Meta `path:"/ws/read" tags:"后台客服面板" method:"post" summary:"消息已读"`
	MsgId  int64 `p:"msg_id" json:"msg_id"`
	Id     int   `p:"id" json:"id"`
}

type ChatMessageReq struct {
	g.Meta `path:"/ws/messages" tags:"后台客服面板" method:"get" summary:"获取消息"`
	Uid    int   `p:"uid" json:"uid" v:"required"`
	Mid    int64 `p:"mid" json:"mid"`
}
type ChatCancelTransferReq struct {
	g.Meta `path:"/ws/transfer/:id/cancel" tags:"后台客服面板" method:"post" summary:"取消转接"`
}

type ChatTransferReq struct {
	g.Meta `path:"/ws/transfer" tags:"后台客服面板" method:"post" summary:"转接用户"`
	UserId int    `v:"required" json:"user_id"`
	ToId   int    `v:"required" json:"to_id"`
	Remark string `v:"max-length:255" json:"remark"`
}

type ChatReqIdReq struct {
	g.Meta `path:"/ws/req-id" tags:"后台客服面板" method:"get" summary:"获取message reqId"`
}

type ChatUserSessionReq struct {
	g.Meta `path:"/ws/sessions/:id" tags:"后台客服面板" method:"get" summary:"获取用户历史session"`
}

type ChatUserListReq struct {
	g.Meta `path:"/ws/chat-users" tags:"后台客服面板" method:"get" summary:"获取客户对应用户列表"`
}

type ChatUserSessionRes []chat.Session

type ChatReqIdRes struct {
	ReqId string `json:"req_id"`
}

type ChatMessageRes []chat.Message

type ChatUserListRes []chat.User

type ChatAcceptRes struct {
	chat.User
}

type SimpleStudent struct {
	Name       string `json:"name"`
	SchoolName string `json:"school_name"`
}

type ChatUserInfoRes struct {
	Phone       string          `json:"phone"`
	RefundCount int             `json:"refund_count"`
	Students    []SimpleStudent `json:"students"`
	OrderCount  int             `json:"order_count"`
}
