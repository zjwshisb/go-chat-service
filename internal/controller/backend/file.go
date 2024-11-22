package backend

import (
	"context"
	"gf-chat/api"
	baseApi "gf-chat/api"
	"gf-chat/internal/library/storage"
	"gf-chat/internal/service"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"strings"
)

var CImage = &cFile{}

type cFile struct {
}

func (c cFile) Store(ctx context.Context, req *api.FileReq) (res *baseApi.NormalRes[*baseApi.File], err error) {
	model, err := storage.Disk().SaveUpload(ctx, req.File, "test")
	if err != nil {
		return nil, err
	}
	file, err := req.File.Open()
	if err != nil {
		return
	}
	defer func() {
		_ = file.Close()
	}()
	mimetype := fileutil.MiMeType(file)
	index := strings.Index(mimetype, "/")
	if index < 0 {
		return nil, gerror.New("unsupported file")
	}
	fileType := mimetype[:index]
	g.Dump(fileType)
	model.CustomerId = service.AdminCtx().GetCustomerId(ctx)
	model.FromId = service.AdminCtx().GetId(ctx)
	model.FromModel = "admin"
	err = service.File().SaveAndFill(ctx, model)
	if err != nil {
		return
	}
	return baseApi.NewResp(service.File().ToApi(model)), nil
}
