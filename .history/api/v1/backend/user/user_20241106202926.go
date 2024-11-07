package user

import (
	"gf-chat/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

type InfoReq struct {
	g.Meta `path:"/me" tags:"管理员" method:"get" summary:"获取管理员信息"`
}

type InfoRes struct {
	Id         uint   `json:"id"`
	CustomerId uint   `json:"customer_id"`
	Username   string `json:"username"`
}

type SettingReq struct {
	g.Meta `path:"/me/settings" tags:"管理员" method:"get" summary:"获取管理员设置"`
}

type UpdateSettingReq struct {
	g.Meta         `path:"/me/settings" tags:"管理员" method:"put" summary:"更新管理员设置"`
	Background     model.ImageFiled `p:"background" json:"background"`
	IsAutoAccept   bool             `p:"is_auto_accept" json:"is_auto_accept"`
	WelcomeContent string           `p:"welcome_content" v:"max-length:512" json:"welcome_content"`
	OfflineContent string           `p:"offline_content" v:"max-length:512" json:"offline_content"`
	Name           string           `p:"name" v:"max-length:20" json:"name"`
	Avatar         model.ImageFiled `p:"avatar" json:"avatar"`
}

type SettingRes struct {
	Background     *model.ImageFiled `json:"background"`
	IsAutoAccept   bool              `json:"is_auto_accept"`
	WelcomeContent string            `json:"welcome_content"`
	OfflineContent string            `json:"offline_content"`
	Name           string            `json:"name"`
	Avatar         *model.ImageFiled `json:"avatar"`
}
