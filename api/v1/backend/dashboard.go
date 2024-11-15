package backend

import (
	"github.com/gogf/gf/v2/frame/g"
)

type DashboardOnlineReq struct {
	g.Meta `path:"/dashboard/online-info" tags:"dashboard" method:"get" summary:"获取在线信息"`
}

type DashboardOnlineUserReq struct {
	g.Meta `path:"/dashboard/online-users" tags:"dashboard" method:"get" summary:"获取在线用户列表"`
}
type DashboardOnlineAdminReq struct {
	g.Meta `path:"/dashboard/online-admins" tags:"dashboard" method:"get" summary:"获取在线客服列表"`
}

type DashboardOnlineUserRes []ChatSimpleUser

type DashboardOnlineRes struct {
	UserCount        uint `json:"user_count"`
	WaitingUserCount uint `json:"waiting_user_count"`
	AdminCount       uint `json:"admin_count"`
}
