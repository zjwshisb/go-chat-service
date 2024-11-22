package backend

import "github.com/gogf/gf/v2/frame/g"

type OptionAutoMessageReq struct {
	g.Meta `path:"/options/auto-messages" tags:"选项" method:"get" summary:"获取快捷回复选项"`
}

type OptionAutoRuleSceneReq struct {
	g.Meta `path:"/options/auto-rule-scenes" tags:"选项" method:"get" summary:"获取回复规则触发场景选项"`
}
