package file

import (
	"context"
	"gf-chat/internal/model"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
)

type StorageAdapter interface {
	Url(path string) string
	Save(ctx context.Context, file *ghttp.UploadFile, relativePath string) (*model.SaveFileOutput, error)
}

type sFile struct {
}

func (s *sFile) Disk(storages ...string) (StorageAdapter, error) {
	def, err := g.Cfg().Get(gctx.New(), "storage.default")
	if err != nil {
		return nil, err
	}
	disk := def.String()
	if len(storages) > 0 {
		disk = storages[0]
	}
	switch disk {
	case "qiniu":
		return qiniu, nil
	}
	return qiniu, nil
}
