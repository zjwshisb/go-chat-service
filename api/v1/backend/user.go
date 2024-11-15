package backend

import (
	"gf-chat/api"
	"github.com/gogf/gf/v2/frame/g"
)

type GetCurrentUserReq struct {
	g.Meta `path:"/user/info" tags:"管理员" method:"get" summary:"获取管理员信息"`
}

type CurrentUserRes struct {
	Id         uint   `json:"id"`
	CustomerId uint   `json:"customer_id"`
	Username   string `json:"username"`
}

type CurrentUserSettingReq struct {
	g.Meta `path:"/user/settings" tags:"管理员" method:"get" summary:"获取管理员设置"`
}

type CurrentUserUpdateSettingReq struct {
	g.Meta         `path:"/user/settings" tags:"管理员" method:"put" summary:"更新管理员设置"`
	Background     api.File `p:"background" json:"background"`
	IsAutoAccept   bool     `p:"is_auto_accept" json:"is_auto_accept"`
	WelcomeContent string   `p:"welcome_content" v:"max-length:512" json:"welcome_content"`
	OfflineContent string   `p:"offline_content" v:"max-length:512" json:"offline_content"`
	Name           string   `p:"name" v:"max-length:20" json:"name"`
	Avatar         api.File `p:"avatar" json:"avatar"`
}

type CurrentUserSettingRes struct {
	Background     api.File `json:"background"`
	IsAutoAccept   bool     `json:"is_auto_accept"`
	WelcomeContent string   `json:"welcome_content"`
	OfflineContent string   `json:"offline_content"`
	Name           string   `json:"name"`
	Avatar         api.File `json:"avatar"`
}

type LoginReq struct {
	g.Meta   `path:"/login" tags:"后台登录" method:"post" summary:"账号密码登录"`
	Username string `v:"required" json:"username"`
	Password string `v:"required" json:"password"`
}

type LoginRes struct {
	Token string `json:"token"`
}
