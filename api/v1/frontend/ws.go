package frontend

import (
	"gf-chat/internal/model/chat"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type ChatConnectReq struct {
	g.Meta `path:"/ws" tags:"C端客服系统" method:"get" summary:"连接websocket服务"`
	Token  string `v:"required" dc:"认证token"`
}

type ChatReqIdReq struct {
	g.Meta `path:"/req-id" tags:"C端客服系统" method:"get" summary:"获取message reqId"`
}

type ChatImageReq struct {
	g.Meta `path:"/image" tags:"C端客服系统" method:"post" summary:"上传图片"`
	File   *ghttp.UploadFile `json:"file" p:"file" type:"file" v:"image" dc:"文件"`
}

type ChatReadReq struct {
	g.Meta `path:"/read" tags:"C端客服系统" method:"post" summary:"消息已读"`
	MsgId  int64 `p:"msg_id"`
}

type ChatRateReq struct {
	g.Meta `path:"/messages/:id/rate" tags:"C端客服系统" method:"post" summary:"消息评分"`
	Rate   int `p:"rate" v:"max:5|min:0"`
}

type ChatImageRes struct {
	Url  string `json:"url"`
	Path string `json:"path"`
}

type ChatReqIdRes struct {
	ReqId string `json:"req_id"`
}

type ChatMessageReq struct {
	g.Meta `path:"/messages" tags:"C端客服系统" method:"get" summary:"获取历史消息"`
	Id     int64 `p:"id"`
	Size   int   `p:"size" d:"20"`
}

type ChatMessageRes []chat.Message
