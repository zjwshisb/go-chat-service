package backend

import (
	"gf-chat/internal/model/chat"

	"github.com/gogf/gf/v2/frame/g"
)

type UserInfoReq struct {
	g.Meta `path:"/ws/chat-user/:id" tags:"后台客服面板" method:"get" summary:"获取用户信息"`
}
type SmsNoticeReq struct {
	g.Meta `path:"/ws/sms-notice" tags:"后台客服面板" method:"post" summary:"发送短信提醒"`
	Uid    uint `json:"uid" v:"required"`
}
type AcceptReq struct {
	g.Meta `path:"/ws/chat-user" tags:"后台客服面板" method:"post" summary:"接入用户"`
	Sid    uint `p:"sid" json:"sid"`
}
type RemoveReq struct {
	g.Meta `path:"/ws/chat-user/:id" tags:"后台客服面板" method:"delete" summary:"移除用户"`
}
type RemoveAllReq struct {
	g.Meta `path:"/ws/chat-user" tags:"后台客服面板" method:"delete" summary:"移除所有失效用户"`
}
type RemoveAllRes struct {
	Ids []uint `json:"ids" d:"移除的用户Id"`
}

type ReadReq struct {
	g.Meta `path:"/ws/read" tags:"后台客服面板" method:"post" summary:"消息已读"`
	MsgId  uint `p:"msg_id" json:"msg_id"`
	Id     uint `p:"id" json:"id"`
}

type MessageReq struct {
	g.Meta `path:"/ws/messages" tags:"后台客服面板" method:"get" summary:"获取消息"`
	Uid    uint `p:"uid" json:"uid" v:"required"`
	Mid    uint `p:"mid" json:"mid"`
}
type CancelTransferReq struct {
	g.Meta `path:"/ws/transfer/:id/cancel" tags:"后台客服面板" method:"post" summary:"取消转接"`
}

type TransferReq struct {
	g.Meta `path:"/ws/transfer" tags:"后台客服面板" method:"post" summary:"转接用户"`
	UserId uint   `v:"required" json:"user_id"`
	ToId   uint   `v:"required" json:"to_id"`
	Remark string `v:"max-length:255" json:"remark"`
}

type ReqIdReq struct {
	g.Meta `path:"/ws/req-id" tags:"后台客服面板" method:"get" summary:"获取message reqId"`
}

type UserSessionReq struct {
	g.Meta `path:"/ws/sessions/:id" tags:"后台客服面板" method:"get" summary:"获取用户历史session"`
}

type UserListReq struct {
	g.Meta `path:"/ws/chat-users" tags:"后台客服面板" method:"get" summary:"获取客户对应用户列表"`
}

type UserSessionRes []chat.Session

type ReqIdRes struct {
	ReqId string `json:"req_id"`
}

type MessageRes []chat.Message

type UserListRes []chat.User

type AcceptRes struct {
	chat.User
}

type SimpleStudent struct {
	Name       string `json:"name"`
	SchoolName string `json:"school_name"`
}

type UserInfoRes struct {
	Phone       string          `json:"phone"`
	RefundCount uint            `json:"refund_count"`
	Students    []SimpleStudent `json:"students"`
	OrderCount  uint            `json:"order_count"`
}
