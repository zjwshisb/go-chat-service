package backend

import (
	"github.com/gogf/gf/v2/frame/g"
)

type CurrentAdminInfoReq struct {
	g.Meta `path:"/current-admin/info" tags:"管理员" method:"get" summary:"获取管理员信息"`
}

type CurrentAdminRes struct {
	Id         uint   `json:"id"`
	CustomerId uint   `json:"customer_id"`
	Username   string `json:"username"`
}

type CurrentAdminSettingReq struct {
	g.Meta `path:"/current-admin/settings" tags:"管理员" method:"get" summary:"获取管理员设置"`
}

type CurrentAdminSetting struct {
	CurrentAdminSettingForm
	AdminId uint
}

type CurrentAdminSettingForm struct {
	Background     *File  `json:"background" v:"api-file:image"`
	IsAutoAccept   bool   `json:"is_auto_accept" v:"boolean"`
	WelcomeContent string `json:"welcome_content" v:"max-length:512"`
	OfflineContent string `json:"offline_content" v:"max-length:512"`
	Name           string `json:"name" v:"max-length:20"`
	Avatar         *File  `json:"avatar" v:"api-file:image"`
}

type CurrentAdminSettingRes struct {
	CurrentAdminSettingForm
}

type CurrentAdminSettingUpdateReq struct {
	g.Meta `path:"/current-admin/settings" tags:"管理员" method:"put" summary:"更新管理员设置"`
	CurrentAdminSettingForm
}

type LoginReq struct {
	g.Meta   `path:"/login" tags:"后台登录" method:"post" summary:"账号密码登录"`
	Username string `v:"required" json:"username"`
	Password string `v:"required" json:"password"`
}

type LoginRes struct {
	Token string `json:"token"`
}
