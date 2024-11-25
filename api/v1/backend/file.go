package backend

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type File struct {
	Id       uint   `json:"id"`
	Path     string `json:"path"`
	Url      string `json:"url"`
	ThumbUrl string `json:"thumb_url"`
	Name     string `json:"name"`
	Type     string `json:"type"`
}

type FileListReq struct {
	g.Meta `path:"/files" tags:"后台文件管理" method:"get" summary:"上传文件"`
	DirId  uint `json:"dir_id"`
}

type FileStoreReq struct {
	g.Meta `path:"/files" tags:"后台图片上传" mine:"multipart/form-data" method:"post" summary:"上传图片"`
	Path   string            `json:"path" p:"path" v:"required" dc:"文件存储路径"`
	File   *ghttp.UploadFile `json:"file" p:"file" type:"file" v:"image" dc:"文件"`
	DirId  uint              `json:"dir_id"`
}

type FileDirStoreReq struct {
	g.Meta `path:"/file-dirs" tags:"后台图片上传" method:"post" summary:"新建文件夹"`
	Pid    uint   `json:"pid"`
	Name   string `json:"name" p:"path" v:"required#请输入文件夹名称" dc:"名称"`
}
