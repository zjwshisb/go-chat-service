package v1

import (
	"gf-chat/api"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type CustomerAdminDetailReq struct {
	g.Meta `path:"/admins/:id" tags:"后台管理员" method:"get" summary:"获取管理员详情"`
	Month  string `json:"month" v:"required" p:"month"`
}

type CustomerAdminListReq struct {
	g.Meta   `path:"/admins" tags:"后台管理员" method:"get" summary:"获取管理员列表"`
	Username string `json:"username"`
	api.Paginate
}

type CustomerAdmin struct {
	Id            uint        `json:"id"`
	Username      string      `json:"username"`
	Avatar        string      `json:"avatar"`
	Online        bool        `json:"online"`
	AcceptedCount uint        `json:"accepted_count"`
	LastOnline    *gtime.Time `json:"last_online"`
	CreatedAt     *gtime.Time `json:"created_at"`
	UpdatedAt     *gtime.Time `json:"updated_at"`
}
