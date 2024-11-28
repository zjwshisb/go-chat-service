package backend

import "github.com/gogf/gf/v2/frame/g"

type OptionAutoMessageReq struct {
	g.Meta `path:"/options/auto-messages" tags:"选项" method:"get" summary:"快捷回复"`
}
type OptionMessageTypeReq struct {
	g.Meta `path:"/options/message-types" tags:"选项" method:"get" summary:"快捷回复类型"`
}

type OptionAutoRuleSceneReq struct {
	g.Meta `path:"/options/auto-rule-scenes" tags:"选项" method:"get" summary:"自动回复规则场景"`
}
type OptionAutoRuleMatchTypeReq struct {
	g.Meta `path:"/options/auto-rule-match-types" tags:"选项" method:"get" summary:"自动回复匹配规则"`
}
type OptionAutoRuleReplyTypeReq struct {
	g.Meta `path:"/options/auto-rule-reply-types" tags:"选项" method:"get" summary:"自动回复回复类型"`
}

type OptionFileTypeReq struct {
	g.Meta `path:"/options/file-types" tags:"选项" method:"get" summary:"文件类型"`
}
