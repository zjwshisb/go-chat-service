package backend

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type ImageReq struct {
	g.Meta `path:"/images" tags:"后台图片上传" mine:"multipart/form-data" method:"post" summary:"上传图片"`
	Path   string            `json:"path" p:"path" v:"required" dc:"文件存储路径"`
	File   *ghttp.UploadFile `json:"file" p:"file" type:"file" v:"image" dc:"文件"`
}

type ImageRes struct {
	Url  string `json:"url"`
	Path string `json:"path"`
}
