package v1

import (
	v1 "gf-chat/api"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type ChatFileListReq struct {
	g.Meta `path:"/chat-files" tags:"聊天文件管理" method:"get" summary:"文件列表"`
	v1.Paginate
}

type ChatFileDeleteReq struct {
	g.Meta `path:"/chat-files/:id" tags:"聊天文件管理" method:"delete" summary:"删除文件"`
}

type ChatFile struct {
	*v1.File
	AdminName string      `json:"admin_name"`
	UserName  string      `json:"user_name"`
	CreatedAt *gtime.Time `json:"created_at"`
}
