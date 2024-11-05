package backend

import "github.com/gogf/gf/v2/frame/g"

type LoginPasswordReq struct {
	g.Meta   `path:"/login/password" tags:"后台登录" method:"post" summary:"账号密码登录"`
	Username string `v:"required" json:"username"`
	Password string `v:"required" json:"password"`
}

type LoginSSoReq struct {
	g.Meta `path:"/login" tags:"后台登录" method:"post" summary:"sso单点登录"`
	Ticket string `p:"ticket" v:"required" json:"ticket"`
}
type LoginRes struct {
	Token string `json:"token"`
}
