package storage

import (
	"context"
	"gf-chat/internal/consts"
	"gf-chat/internal/model/entity"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
)

type Adapter interface {
	Url(path string) string
	ThumbUrl(path string) string
	SaveUpload(ctx context.Context, file *ghttp.UploadFile, relativePath string) (*entity.CustomerChatFiles, error)
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
