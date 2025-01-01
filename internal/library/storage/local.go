package storage

import (
	"context"
	"gf-chat/internal/consts"
	"gf-chat/internal/model"
	"gf-chat/internal/model/entity"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"os"
	"strings"
)

var pathSeparator = "/"

func newLocal() *localAdapter {
	ctx := gctx.New()
	serverRootVar, err := g.Cfg().Get(ctx, "server.serverRoot")
	if err != nil {
		g.Log().Errorf(ctx, "%+v", err)
		return nil
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
	ctx := gctx.New()
	host, err := g.Config().Get(ctx, "app.host")
	if err != nil {
		g.Log().Error(ctx)
		return ""
	}
	return host.String() + pathSeparator + path
}
func (s *localAdapter) ThumbUrl(path string) string {
	return s.Url(path)
}
func (s *localAdapter) Delete(path string) error {
	dir := s.serverRoot + pathSeparator + path
	return os.Remove(dir)
}

func (s *localAdapter) SaveUpload(_ context.Context, file *ghttp.UploadFile, relativePath string) (files *model.CustomerChatFile, err error) {
	var fullPath string
	if strings.HasSuffix(relativePath, pathSeparator) {
		relativePath = relativePath[:len(relativePath)-1]
	}
	fullPath = s.serverRoot + pathSeparator + relativePath
	name, err := file.Save(fullPath, true)
	if err != nil {
		return
	}
	if relativePath != "" {
		relativePath = relativePath + pathSeparator
	}
	relativePath += name
	return &model.CustomerChatFile{
		CustomerChatFiles: entity.CustomerChatFiles{
			Name: file.Filename,
			Path: relativePath,
			Disk: consts.StorageLocal,
		},
	}, nil
}
