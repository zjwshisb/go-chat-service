package backend

import (
	"github.com/gogf/gf/v2/frame/g"
)

type ChatConnectReq struct {
	g.Meta `path:"/ws" tags:"后台websocket链接" method:"get" summary:"连接websocket服务"`
	Token  string `v:"required" dc:"认证token"`
}
