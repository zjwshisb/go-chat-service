package frontend

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type SettingReq struct {
	g.Meta `path:"/setting" tags:"C端客服系统" method:"get" summary:"获取设置"`
}

type SettingRes struct {
	IsShowQueue bool `json:"is_show_queue"`
	IsShowRead  bool `json:"is_show_read"`
}

type ChatConnectReq struct {
	g.Meta `path:"/ws" tags:"C端客服系统" method:"get" summary:"连接websocket服务"`
	Token  string `v:"required" dc:"认证token"`
}

type ChatReqIdReq struct {
	g.Meta `path:"/req-id" tags:"C端客服系统" method:"get" summary:"获取message reqId"`
}

type ChatReadReq struct {
	g.Meta `path:"/read" tags:"C端客服系统" method:"post" summary:"消息已读"`
	MsgId  uint `p:"msg_id"`
}

type ChatRateReq struct {
	g.Meta `path:"/messages/:id/rate" tags:"C端客服系统" method:"post" summary:"消息评分"`
	Rate   uint `p:"rate" v:"max:5|min:0"`
}

type FileStoreReq struct {
	g.Meta `path:"/files" tags:"上传图片" mine:"multipart/form-data" method:"post" summary:"上传文件"`
	File   *ghttp.UploadFile `json:"file" p:"file" type:"file" v:"required|file" dc:"文件"`
}
type ChatReqId struct {
	ReqId string `json:"req_id"`
}

type ChatMessageReq struct {
	g.Meta   `path:"/messages" tags:"C端客服系统" method:"get" summary:"获取历史消息"`
	Id       uint `p:"id"`
	PageSize int  `d:"20"`
}
