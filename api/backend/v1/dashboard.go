package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type DashboardOnlineReq struct {
	g.Meta `path:"/dashboard/online-info" tags:"dashboard" method:"get" summary:"获取在线信息"`
}
type DashboardWaitingUserInfoReq struct {
	g.Meta `path:"/dashboard/waiting-user-info" tags:"dashboard" method:"get" summary:"获取等待用户列表"`
}
type DashboardOnlineUserInfoReq struct {
	g.Meta `path:"/dashboard/online-user-info" tags:"dashboard" method:"get" summary:"获取在线用户列表"`
}
type DashboardAdminInfoReq struct {
	g.Meta `path:"/dashboard/admin-info" tags:"dashboard" method:"get" summary:"获取在线客服列表"`
}

type DashboardAdminInfo struct {
	Admins []ChatSimpleUser `json:"admins"`
	Total  int              `json:"total"`
}

type DashboardOnlineUserInfo struct {
	Users       []ChatSimpleUser `json:"users"`
	ActiveCount int              `json:"active_count"`
}

type DashboardWaitingUserInfo struct {
	Users      []ChatSimpleUser `json:"users"`
	TodayTotal int              `json:"today_total"`
}
