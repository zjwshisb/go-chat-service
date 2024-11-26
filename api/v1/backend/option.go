package backend

import "github.com/gogf/gf/v2/frame/g"

type OptionAutoMessageReq struct {
	g.Meta `path:"/options/auto-messages" tags:"选项" method:"get" summary:"选项"`
}

type OptionAutoRuleSceneReq struct {
	g.Meta `path:"/options/auto-rule-scenes" tags:"选项" method:"get" summary:"选项"`
}

type OptionFileTypeReq struct {
	g.Meta `path:"/options/file-types" tags:"选项" method:"get" summary:"选项"`
}
