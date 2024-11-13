package autorule

import (
	api "gf-chat/api/v1/backend/automessage"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type ListReq struct {
	g.Meta   `path:"/auto-rules" tags:"后台自动回复规则" method:"get" summary:"获取自动回复规则"`
	PageSize int `d:"20" json:"pageSize" v:"max:100"`
	Current  int `d:"1" dc:"页码" json:"current"`
}
type StoreReq struct {
	g.Meta    `path:"/auto-rules" tags:"后台自动回复规则" method:"post" summary:"新增自动回复规则"`
	Name      string   `json:"name" p:"name" v:"required|max-length:32|unique:customer_chat_auto_rules,name#||已存在相同名字的规则"`
	Match     string   `json:"match" p:"match" v:"required"`
	MatchType string   `json:"match_type" p:"match_type" v:"required"`
	ReplyType string   `json:"reply_type" p:"reply_type" v:"required"`
	MessageId uint     `json:"message_id" p:"message_id"`
	IsOpen    bool     `json:"is_open" p:"is_open" v:"boolean"`
	Sort      uint     `json:"sort" p:"sort" v:"required|max:128|min:0"`
	Scenes    []string `json:"scenes" p:"scenes" `
}
type UpdateReq struct {
	g.Meta    `path:"/auto-rules/:id" tags:"后台自动回复规则" method:"put" summary:"编辑自动回复规则"`
	Name      string   `json:"name" p:"name" v:"required|max-length:32|unique:customer_chat_auto_rules,name,id#||已存在相同名字的规则"`
	Match     string   `json:"match" p:"match" v:"required"`
	MatchType string   `json:"match_type" p:"match_type" v:"required"`
	ReplyType string   `json:"reply_type" p:"reply_type" v:"required"`
	MessageId uint     `json:"message_id" p:"message_id"`
	IsOpen    bool     `json:"is_open" p:"is_open" v:"boolean"`
	Sort      uint     `json:"sort" p:"sort" v:"required|max:128|min:0"`
	Scenes    []string `json:"scenes" p:"scenes" `
}

type DeleteReq struct {
	g.Meta `path:"/auto-rules/:id" tags:"后台自动回复规则" method:"delete" summary:"删除自动回复规则"`
}

type FormReq struct {
	g.Meta `path:"/auto-rules/:id/form" tags:"后台自动回复规则" method:"get" summary:"获取自动回复规则表单"`
}

type FormRes struct {
	IsOpen    bool     `json:"is_open"`
	Match     string   `json:"match"`
	MatchType string   `json:"match_type"`
	MessageId uint     `json:"message_id"`
	Name      string   `json:"name"`
	ReplyType string   `json:"reply_type"`
	Scenes    []string `json:"scenes"`
	Sort      uint     `json:"sort"`
}

type ListItem struct {
	Id         uint         `json:"id"`
	Name       string       `json:"name"`
	Match      string       `json:"match"`
	MatchType  string       `json:"match_type"`
	ReplyType  string       `json:"reply_type"`
	MessageId  uint         `json:"message_id"`
	Key        string       `gorm:"key" json:"key"`
	Sort       uint         `json:"sort"`
	IsOpen     bool         `json:"is_open"`
	Count      uint         `json:"count"`
	CreatedAt  *gtime.Time  `json:"created_at"`
	UpdatedAt  *gtime.Time  `json:"updated_at"`
	EventLabel string       `json:"event_label"`
	Scenes     []string     `json:"scenes"`
	Message    api.ListItem `json:"message"`
}

type ListRes struct {
	Items []ListItem `json:"items"`
	Total int        `json:"total"`
}
