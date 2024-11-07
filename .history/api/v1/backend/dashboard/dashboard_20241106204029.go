package dashboard

import (
	"gf-chat/internal/model/chat"

	"github.com/gogf/gf/v2/frame/g"
)

type OnlineReq struct {
	g.Meta `path:"/dashboard/online-info" tags:"dashboard" method:"get" summary:"获取在线信息"`
}

type OnlineUserReq struct {
	g.Meta `path:"/dashboard/online-users" tags:"dashboard" method:"get" summary:"获取在线用户列表"`
}
type OnlineAdminReq struct {
	g.Meta `path:"/dashboard/online-admins" tags:"dashboard" method:"get" summary:"获取在线客服列表"`
}

type OnlineUserRes []chat.SimpleUser

type OnlineRes struct {
	UserCount        uint `json:"user_count"`
	WaitingUserCount uint `json:"waiting_user_count"`
	AdminCount       uint `json:"admin_count"`
}
