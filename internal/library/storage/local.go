package storage

import (
	"context"
	"gf-chat/internal/consts"
	"gf-chat/internal/model/entity"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"log"
	"os"
	"strings"
)

var pathSeparator = string(os.PathSeparator)

func newLocal() *localAdapter {
	ctx := gctx.New()
	serverRootVar, err := g.Cfg().Get(ctx, "server.serverRoot")
	if err != nil {
		log.Fatal(err)
	}
	serverRoot := serverRootVar.String()
	if strings.HasSuffix(serverRoot, pathSeparator) {
		serverRoot = serverRoot[:len(serverRoot)-1]
	}
	return &localAdapter{
		serverRoot: serverRoot,
	}
}

type localAdapter struct {
	serverRoot string
}

func (s *localAdapter) Url(path string) string {
	return path
}
func (s *localAdapter) ThumbUrl(path string) string {
	return s.Url(path)
}

func (s *localAdapter) SaveUpload(ctx context.Context, file *ghttp.UploadFile, relativePath string) (files *entity.CustomerChatFiles, err error) {
	var fullPath string
	if strings.HasPrefix(relativePath, pathSeparator) {
		fullPath = s.serverRoot + relativePath
	} else {
		fullPath = s.serverRoot + pathSeparator + relativePath
	}
	name, err := file.Save(fullPath, true)
	if err != nil {
		return
	}
	return &entity.CustomerChatFiles{
		Path: fullPath + pathSeparator + name,
		Disk: consts.StorageLocal,
	}, nil
}
