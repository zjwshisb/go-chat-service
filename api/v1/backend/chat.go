package backend

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type GetUserChatInfoReq struct {
	g.Meta `path:"/ws/chat-user/:id" tags:"后台客服面板" method:"get" summary:"获取用户信息"`
}

type AcceptUserReq struct {
	g.Meta `path:"/ws/chat-user" tags:"后台客服面板" method:"post" summary:"接入用户"`
	Sid    uint `p:"sid" json:"sid"`
}
type RemoveUserReq struct {
	g.Meta `path:"/ws/chat-user/:id" tags:"后台客服面板" method:"delete" summary:"移除用户"`
}
type RemoveAllUserReq struct {
	g.Meta `path:"/ws/chat-user" tags:"后台客服面板" method:"delete" summary:"移除所有失效用户"`
}
type RemoveAllUserRes struct {
	Ids []uint `json:"ids" d:"移除的用户Id"`
}

type MessageReadReq struct {
	g.Meta `path:"/ws/read" tags:"后台客服面板" method:"post" summary:"消息已读"`
	MsgId  uint `p:"msg_id" json:"msg_id"`
	Id     uint `p:"id" json:"id"`
}

type GetMessageReq struct {
	g.Meta `path:"/ws/messages" tags:"后台客服面板" method:"get" summary:"获取消息"`
	Uid    uint `p:"uid" json:"uid" v:"required"`
	Mid    uint `p:"mid" json:"mid"`
}
type CancelTransferReq struct {
	g.Meta `path:"/ws/transfer/:id/cancel" tags:"后台客服面板" method:"post" summary:"取消转接"`
}

type StoreTransferReq struct {
	g.Meta `path:"/ws/transfer" tags:"后台客服面板" method:"post" summary:"转接用户"`
	UserId uint   `v:"required" json:"user_id"`
	ToId   uint   `v:"required" json:"to_id"`
	Remark string `v:"max-length:255" json:"remark"`
}

type ReqIdReq struct {
	g.Meta `path:"/ws/req-id" tags:"后台客服面板" method:"get" summary:"获取message reqId"`
}

type GetUserSessionReq struct {
	g.Meta `path:"/ws/sessions/:id" tags:"后台客服面板" method:"get" summary:"获取用户历史session"`
}

type UserListReq struct {
	g.Meta `path:"/ws/chat-users" tags:"后台客服面板" method:"get" summary:"获取客户对应用户列表"`
}

type ReqIdRes struct {
	ReqId string `json:"req_id"`
}

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

type ChatUser struct {
	Id           uint          `json:"id"`
	Username     string        `json:"username"`
	LastChatTime *gtime.Time   `json:"last_chat_time"`
	Disabled     bool          `json:"disabled"`
	Online       bool          `json:"online"`
	LastMessage  *ChatMessage  `json:"last_message"`
	Unread       uint          `json:"unread"`
	Avatar       string        `json:"avatar"`
	Platform     string        `json:"platform"`
	Messages     []ChatMessage `json:"messages"`
}

type ChatSimpleUser struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
}

type MessageRes []*ChatMessage

type UserListRes []ChatUser

type AcceptRes struct {
	User ChatUser
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
type UserInfoItem struct {
	Name        string
	Label       string
	Description string
}
type UserInfoRes struct {
	Phone string `json:"phone"`
}
