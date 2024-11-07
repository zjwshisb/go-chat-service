package systemrule

import "github.com/gogf/gf/v2/frame/g"

type ListReq struct {
	g.Meta `path:"/system-auto-rules" tags:"系统规则" method:"get" summary:"获取系统规则设置"`
}

type UpdateReq struct {
	g.Meta `path:"/system-auto-rules" tags:"系统规则" method:"put" summary:"更新系统规则设置"`
	Data   map[string]string `p:"data" json:"data"`
}

type ListRes []ListItem

type ListItem struct {
	Id        uint   `json:"id"`
	MessageId uint   `json:"message_id"`
	Name      string `json:"name"`
}
