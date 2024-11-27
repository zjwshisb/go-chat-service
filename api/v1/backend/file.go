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
	g.Meta `path:"/files" tags:"后台文件管理" method:"get" summary:"文件列表"`
	DirId  uint `json:"dir_id"`
}

type FileStoreReq struct {
	g.Meta `path:"/files" tags:"后台文件管理" mine:"multipart/form-data" method:"post" summary:"上传文件"`
	File   *ghttp.UploadFile `json:"file" p:"file" type:"file" v:"required|file" dc:"文件"`
	Pid    uint              `json:"pid" v:"file-dir"`
}

type FileDirStoreReq struct {
	g.Meta `path:"/file-dirs" tags:"后台文件管理" method:"post" summary:"新建目录"`
	Pid    uint   `json:"pid" v:"file-dir"`
	Name   string `json:"name" p:"path" v:"required|length:1,20#请输入文件夹名称|名称最长20个字符" dc:"名称"`
}

type FileUpdateReq struct {
	g.Meta `path:"/files/:id" tags:"后台文件管理" method:"put" summary:"修改文件名"`
	Name   string `json:"name" p:"path" v:"required|length:1,20#请输入文件夹名称|名称最长20个字符" dc:"名称"`
}

type FileDeleteReq struct {
	g.Meta `path:"/files/:id" tags:"后台文件管理" method:"delete" summary:"删除文件"`
}
