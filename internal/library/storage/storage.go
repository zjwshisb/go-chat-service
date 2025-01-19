package storage

import (
	"context"
	"gf-chat/internal/consts"
	"gf-chat/internal/model"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"mime/multipart"
	"net/http"
	"strings"
)

var DefaultFileSize = map[string]int{
	consts.FileTypeImage: 5 * 1024 * 1024,
	consts.FileTypeVideo: 10 * 1024 * 1024,
	consts.FileTypeAudio: 5 * 1024 * 1024,
	consts.FileTypePdf:   20 * 1024 * 1024,
}

var (
	local = newLocal()
	qiniu = newQiniu()
)

type Adapter interface {
	Url(path string) string
	ThumbUrl(path string) string
	SaveUpload(ctx context.Context, file *ghttp.UploadFile, relativePath string) (*model.CustomerChatFile, error)
	Delete(path string) error
}

func FileType(file multipart.File) (string, error) {
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return "", err
	}
	mimetype := http.DetectContentType(buffer)
	g.Dump(mimetype)
	index := strings.Index(mimetype, "/")
	if index < 0 {
		return "", gerror.NewCode(gcode.CodeValidationFailed, "不支持的文件类型")
	}
	types := mimetype[:index]
	if types == "application" {
		types = mimetype[index+1:]
	}
	if _, exist := DefaultFileSize[types]; !exist {
		return "", gerror.NewCode(gcode.CodeValidationFailed, "不支持的文件类型")
	}
	return types, nil
}

func Disk(storages ...string) Adapter {
	def, err := g.Cfg().Get(gctx.New(), "storage.default")
	if err != nil {
		panic(err)
	}
	disk := def.String()
	if len(storages) > 0 {
		disk = storages[0]
	}
	var adapter Adapter
	switch disk {
	case consts.StorageQiniu:
		adapter = qiniu
		break
	case consts.StorageLocal:
		adapter = local
		break
	default:
		panic("暂不支持该存储引擎:" + disk)
	}
	return adapter
}
