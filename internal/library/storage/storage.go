package storage

import (
	"context"
	"gf-chat/internal/consts"
	"gf-chat/internal/model"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"mime/multipart"
	"strings"
)

var DefaultFileSize = map[string]int{
	consts.FileTypeImage: 5 * 1024 * 1024,
	consts.FileTypeVideo: 10 * 1024 * 1024,
	consts.FileTypeAudio: 5 * 1024 * 1024,
}

type Adapter interface {
	Url(path string) string
	ThumbUrl(path string) string
	SaveUpload(ctx context.Context, file *ghttp.UploadFile, relativePath string) (*model.CustomerChatFile, error)
}

func FileType(file multipart.File) (string, error) {
	mimetype := fileutil.MiMeType(file)
	index := strings.Index(mimetype, "/")
	if index < 0 {
		return "", gerror.New("unsupported file type")
	}
	types := mimetype[:index]
	if _, exist := DefaultFileSize[types]; !exist {
		return "", gerror.New("unsupported file type")
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
		adapter = newQiniu()
		break
	case consts.StorageLocal:
		adapter = newLocal()
		break
	default:
		panic("暂不支持该存储引擎:" + disk)
	}
	return adapter
}
