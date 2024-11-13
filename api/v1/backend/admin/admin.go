package admin

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type DetailReq struct {
	g.Meta `path:"/admins/:id" tags:"后台管理员" method:"get" summary:"获取管理员详情"`
	Month  string `json:"month" v:"required" p:"month"`
}

type DetailRes struct {
}

type ListReq struct {
	g.Meta   `path:"/admins" tags:"后台管理员" method:"get" summary:"获取管理员列表"`
	PageSize int `d:"20" json:"pageSize" v:"max:100"`
	Current  int `d:"1" dc:"页码" json:"current"`
}

type StoreReq struct {
	g.Meta   `path:"/admins" tags:"后台管理员" method:"get" summary:"获取管理员列表"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ListItem struct {
	Id            uint        `json:"id"`
	Username      string      `json:"username"`
	Avatar        string      `json:"avatar"`
	Online        bool        `json:"online"`
	AcceptedCount uint        `json:"accepted_count"`
	LastOnline    *gtime.Time `json:"last_online"`
}
