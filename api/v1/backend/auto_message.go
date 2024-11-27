package backend

import (
	"gf-chat/api"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/gogf/gf/v2/frame/g"
)

type AutoMessageListReq struct {
	g.Meta `path:"/auto-messages" tags:"后台快捷回复" method:"get" summary:"获取快捷回复详情"`
	api.Paginate
	Type string `json:"type"`
	Name string `json:"name"`
}

type AutoMessageFormReq struct {
	g.Meta `path:"/auto-messages/:id/form" tags:"后台快捷回复" method:"get" summary:"获取编辑表单数据"`
}

type AutoMessageFormRes struct {
	AutoMessageForm
}

type AutoMessageNavigator struct {
	Url   string `json:"url" v:"required-if:type,navigator|max-length:512"`
	Title string `json:"title" v:"required-if:type,navigator|max-length:32"`
	Image *File  `json:"image" v:"required-if:type,navigator|api-file:image"`
}

type AutoMessageForm struct {
	Type      string                `json:"type" v:"required|auto-message-type"`
	Name      string                `json:"name" v:"required|max-length:32|unique:customer_chat_auto_messages,name#||已存在相同名字的消息"`
	Content   string                `json:"content" v:"required-if:type,text|max-length:512"`
	Navigator *AutoMessageNavigator `json:"navigator" v:"required-if:type,navigator"`
	File      *File                 `json:"file" v:"required-if:type,image|api-file"`
}

type AutoMessageStoreReq struct {
	g.Meta `path:"/auto-messages" tags:"后台快捷回复" method:"post" summary:"新增快捷回复"`
	AutoMessageForm
}

type AutoMessageUpdateReq struct {
	g.Meta `path:"/auto-messages/:id" tags:"后台快捷回复" method:"put" summary:"修改快捷回复"`
	AutoMessageForm
}
type AutoMessageDeleteReq struct {
	g.Meta `path:"/auto-messages/:id" tags:"后台快捷回复" method:"delete" summary:"删除快捷回复"`
}

type AutoMessageListItem struct {
	Id         uint                  `json:"id"`
	Name       string                `json:"name"`
	Type       string                `json:"type"`
	Content    string                `json:"content"`
	File       *File                 `json:"file"`
	Navigator  *AutoMessageNavigator `json:"navigator"`
	CreatedAt  *gtime.Time           `json:"created_at"`
	UpdatedAt  *gtime.Time           `json:"updated_at"`
	RulesCount uint                  `json:"rules_count"`
}
