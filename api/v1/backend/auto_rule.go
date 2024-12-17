package backend

import (
	"gf-chat/api/v1"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type AutoRuleListReq struct {
	g.Meta `path:"/auto-rules" tags:"后台自动回复规则" method:"get" summary:"获取自动回复规则"`
	v1.Paginate
	ReplyType string `json:"reply_type"`
	Name      string `json:"name"`
	MatchType string `json:"match_type"`
	IsOpen    *bool  `json:"is_open"`
}

type AutoRuleForm struct {
	Name      string   `json:"name"  v:"required|max-length:32|unique:customer_chat_auto_rules,name#||已存在相同名字的规则"`
	Match     string   `json:"match"  v:"required"`
	MatchType string   `json:"match_type"  v:"required|auto-rule-match-type"`
	ReplyType string   `json:"reply_type"  v:"required|auto-rule-reply-type"`
	MessageId uint     `json:"message_id" v:"required-if:reply_type,message|exists:customer_chat_auto_messages"`
	IsOpen    bool     `json:"is_open" p:"is_open" v:"boolean"`
	Sort      uint     `json:"sort" p:"sort" v:"required|max:10000|min:0"`
	Scenes    []string `json:"scenes" p:"scenes" v:"required-if:reply_type,message|foreach|auto-rule-scene"`
}

type AutoRuleStoreReq struct {
	g.Meta `path:"/auto-rules" tags:"后台自动回复规则" method:"post" summary:"新增自动回复规则"`
	AutoRuleForm
}
type AutoRuleUpdateReq struct {
	g.Meta `path:"/auto-rules/:id" tags:"后台自动回复规则" method:"put" summary:"编辑自动回复规则"`
	AutoRuleForm
}

type AutoRuleDeleteReq struct {
	g.Meta `path:"/auto-rules/:id" tags:"后台自动回复规则" method:"delete" summary:"删除自动回复规则"`
}

type AutoRuleFormReq struct {
	g.Meta `path:"/auto-rules/:id/form" tags:"后台自动回复规则" method:"get" summary:"获取自动回复规则表单"`
}

type AutoRule struct {
	Id         uint         `json:"id"`
	Name       string       `json:"name"`
	Match      string       `json:"match"`
	MatchType  string       `json:"match_type"`
	ReplyType  string       `json:"reply_type"`
	Sort       uint         `json:"sort"`
	IsOpen     bool         `json:"is_open"`
	Count      uint         `json:"count"`
	CreatedAt  *gtime.Time  `json:"created_at"`
	UpdatedAt  *gtime.Time  `json:"updated_at"`
	EventLabel string       `json:"event_label"`
	Scenes     []string     `json:"scenes"`
	Message    *AutoMessage `json:"message"`
}
