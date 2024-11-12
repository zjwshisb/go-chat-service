package automessage

import (
	"gf-chat/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

type ListReq struct {
	g.Meta   `path:"/auto-messages" tags:"后台快捷回复" method:"get" summary:"获取快捷回复详情"`
	PageSize int    `d:"20" json:"pageSize" v:"max:100"`
	Current  int    `d:"1" dc:"页码" json:"current"`
	Type     string `p:"type"`
}

type FormReq struct {
	g.Meta `path:"/auto-messages/:id/form" tags:"后台快捷回复" method:"get" summary:"获取编辑表单数据"`
}

type StoreReq struct {
	g.Meta  `path:"/auto-messages" tags:"后台快捷回复" method:"post" summary:"新增快捷回复"`
	Name    string `json:"name" p:"name" v:"required|max-length:32|unique:customer_chat_auto_messages,name#||已存在相同名字的消息"`
	Type    string `json:"type" p:"type" v:"required|auto-message-type"`
	Content string `json:"content" p:"content"  v:"required|max-length:512"`
	Title   string `json:"title" p:"title"  v:"max-length:32"`
	Url     string `json:"url" p:"url" v:"max-length:512"`
}
type UpdateReq struct {
	g.Meta  `path:"/auto-messages/:id" tags:"后台快捷回复" method:"put" summary:"修改快捷回复"`
	Name    string `json:"name" p:"name" v:"required|max-length:32|unique:customer_chat_auto_messages,name,id#||已存在相同名字的消息"`
	Content string `json:"content" p:"content"  v:"required|max-length:512"`
	Title   string `json:"title" p:"title"  v:"max-length:32"`
	Url     string `json:"url" p:"url" v:"max-length:512"`
}
type DeleteReq struct {
	g.Meta `path:"/auto-messages/:id" tags:"后台快捷回复" method:"delete" summary:"删除快捷回复"`
}

type OptionReq struct {
	g.Meta `path:"/options/auto-messages" tags:"后台快捷回复" method:"get" summary:"获取快捷回复选项"`
}

type ListRes struct {
	Total int                         `json:"total"`
	Items []model.AutoMessageListItem `json:"items"`
}

type FormRes struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content any    `json:"content"`
	Title   string `json:"title"`
	Url     string `json:"url"`
}
