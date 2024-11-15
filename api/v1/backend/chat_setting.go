package backend

import (
	"gf-chat/api"
	"github.com/gogf/gf/v2/frame/g"
)

type ChatSettingListReq struct {
	g.Meta `path:"/settings" tags:"后台系统设置" method:"get" summary:"获取系统设置列表"`
}
type ChatSettingUpdateReq struct {
	g.Meta `path:"/settings/:id" tags:"后台系统设置" method:"put" summary:"修改系统设置列表"`
	Value  string `p:"value" v:"required" json:"value"`
}

type ChatSettingListItem struct {
	Id          uint         `json:"id"`
	Name        string       `json:"name"`
	Value       any          `json:"value"`
	Options     []api.Option `json:"options"`
	Title       string       `json:"title"`
	Type        string       `json:"type"`
	Description string       `json:"description"`
}

type ListRes []ChatSettingListItem
