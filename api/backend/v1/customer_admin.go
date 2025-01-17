package v1

import (
	"gf-chat/api"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type CustomerAdminListReq struct {
	g.Meta   `path:"/admins" tags:"后台管理员" method:"get" summary:"获取管理员列表"`
	Username string `json:"username"`
	api.Paginate
}

type StoreCustomerAdminReq struct {
	g.Meta `path:"/admins" tags:"后台管理员" method:"post" summary:"新增管理员"`
	CustomerAdminForm
}

type UpdateCustomerAdminReq struct {
	g.Meta `path:"/admins/:id" tags:"后台管理员" method:"put" summary:"修改管理员"`
}

type CustomerAdminForm struct {
	Username string `json:"username" v:"required|max-length:32|unique:customer_admins,username#||登录账号不可用"`
	Password string `json:"password" v:"required|max-length:32"`
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
