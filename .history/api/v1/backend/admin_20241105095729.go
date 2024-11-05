package backend

import (
	"gf-chat/internal/model"
	"github.com/gogf/gf/v2/frame/g"
	"time"
)

type AdminDetailReq struct {
	g.Meta `path:"/admins/:id" tags:"后台管理员" method:"get" summary:"获取管理员详情"`
	Month  string `json:"month" v:"required" p:"month"`
}

type AdminDetailRes struct {
	Chart []*model.ChartLine     `json:"chart"`
	Admin *model.AdminDetailInfo `json:"admin"`
}

type AdminIndexReq struct {
	g.Meta   `path:"/admins" tags:"后台管理员" method:"get" summary:"获取管理员列表"`
	PageSize int `d:"20" json:"pageSize" v:"max:100"`
	Current  int `d:"1" dc:"页码" json:"current"`
}

type AdminListItem struct {
	Id            uint       `json:"id"`
	Username      string     `json:"username"`
	Avatar        string     `json:"avatar"`
	Online        bool       `json:"online"`
	AcceptedCount int        `json:"accepted_count"`
	LastOnline    *time.Time `json:"last_online"`
}

type AdminIndexRes struct {
	Total int             `json:"total"`
	Items []AdminListItem `json:"items"`
}
