package backend

import (
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/gogf/gf/v2/frame/g"
)

type AutoMessageListReq struct {
	g.Meta   `path:"/auto-messages" tags:"后台快捷回复" method:"get" summary:"获取快捷回复详情"`
	PageSize int    `d:"20" json:"pageSize" v:"max:100"`
	Current  int    `d:"1" dc:"页码" json:"current"`
	Type     string `p:"type"`
}

type AutoMessageFormReq struct {
	g.Meta `path:"/auto-messages/:id/form" tags:"后台快捷回复" method:"get" summary:"获取编辑表单数据"`
}

type AutoMessageStoreReq struct {
	g.Meta  `path:"/auto-messages" tags:"后台快捷回复" method:"post" summary:"新增快捷回复"`
	Name    string `json:"name" p:"name" v:"required|max-length:32|unique:customer_chat_auto_messages,name#||已存在相同名字的消息"`
	Type    string `json:"type" p:"type" v:"required|auto-message-type"`
	Content string `json:"content" p:"content"  v:"required|max-length:512"`
	Title   string `json:"title" p:"title"  v:"max-length:32"`
	Url     string `json:"url" p:"url" v:"max-length:512"`
}
type AutoMessageUpdateReq struct {
	g.Meta  `path:"/auto-messages/:id" tags:"后台快捷回复" method:"put" summary:"修改快捷回复"`
	Name    string `json:"name" p:"name" v:"required|max-length:32|unique:customer_chat_auto_messages,name,id#||已存在相同名字的消息"`
	Content string `json:"content" p:"content"  v:"required|max-length:512"`
	Title   string `json:"title" p:"title"  v:"max-length:32"`
	Url     string `json:"url" p:"url" v:"max-length:512"`
}
type AutoMessageDeleteReq struct {
	g.Meta `path:"/auto-messages/:id" tags:"后台快捷回复" method:"delete" summary:"删除快捷回复"`
}

type AutoMessageOptionReq struct {
	g.Meta `path:"/options/auto-messages" tags:"后台快捷回复" method:"get" summary:"获取快捷回复选项"`
}

type AutoMessageListItem struct {
	Id         uint        `json:"id"`
	Name       string      `json:"name"`
	Type       string      `json:"type"`
	Content    string      `json:"content"`
	Url        string      `json:"url"`
	Title      string      `json:"title"`
	CreatedAt  *gtime.Time `json:"created_at"`
	UpdatedAt  *gtime.Time `json:"updated_at"`
	RulesCount uint        `json:"rules_count"`
}

type AutoMessageFormRes struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content any    `json:"content"`
	Title   string `json:"title"`
	Url     string `json:"url"`
}
