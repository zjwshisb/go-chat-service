package v1

import (
	"gf-chat/api"
	"github.com/gogf/gf/v2/frame/g"
)

type TransferListReq struct {
	g.Meta `path:"/transfers" tags:"后台转接记录" method:"get" summary:"获取转接记录列表"`
	api.Paginate
	Username      string `dc:"用户名"`
	FromAdminName string `dc:"转接客服名称"`
	ToAdminName   string `dc:"转接给客服名称"`
}

type TransferCancelReq struct {
	g.Meta `path:"/transfers/:id/cancel" tags:"后台转接记录" method:"post" summary:"取消转接记录"`
}

type Transfer struct {
	Id            int64  `json:"id"`
	SessionId     uint64 `json:"session_id"`
	UserId        int    `json:"user_id"`
	Remark        string `json:"remark"`
	FromAdminName string `json:"from_admin_name"`
	ToAdminName   string `json:"to_admin_name"`
	Username      string `json:"username"`
	CreatedAt     int64  `json:"created_at"`
	AcceptedAt    int64  `json:"accepted_at"`
	CanceledAt    int64  `json:"canceled_at"`
}
