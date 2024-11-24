package backend

import (
	"context"
	baseApi "gf-chat/api"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/library/storage"
	"gf-chat/internal/model"
	"gf-chat/internal/service"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/frame/g"
)

var CImage = &cFile{}

type cFile struct {
}

func (c cFile) Index(ctx context.Context, req *api.FileListReq) (res *baseApi.ListRes[*api.File], err error) {
	where := g.Map{
		"customer_id": service.AdminCtx().GetCustomerId(ctx),
	}
	if req.Type != "" {
		where["type"] = req.Type
	}
	if req.DirId != 0 {
		where["dir_id"] = req.DirId
	}
	if req.LastId != 0 {
		where["id < ?"] = req.LastId
	}
	files, err := service.File().All(ctx, where, nil, "id desc", 50)
	if err != nil {
		return
	}
	count, err := service.File().Count(ctx, where)
	if err != nil {
		return
	}
	apiFiles := slice.Map(files, func(index int, item *model.CustomerChatFile) *api.File {
		return service.File().ToApi(item)
	})

	return baseApi.NewListResp(apiFiles, count), nil
}

func (c cFile) Store(ctx context.Context, req *api.FileStoreReq) (res *baseApi.NormalRes[*api.File], err error) {
	file, err := req.File.Open()
	if err != nil {
		return
	}
	defer func() {
		_ = file.Close()
	}()
	fileType, err := storage.FileType(file)
	if err != nil {
		return
	}
	fileModel, err := storage.Disk().SaveUpload(ctx, req.File, "")
	if err != nil {
		return nil, err
	}
	fileModel.Type = fileType
	fileModel.CustomerId = service.AdminCtx().GetCustomerId(ctx)
	fileModel.FromId = service.AdminCtx().GetId(ctx)
	fileModel.FromModel = "admin"
	err = service.File().SaveAndFill(ctx, fileModel)
	if err != nil {
		return
	}
	return baseApi.NewResp(service.File().ToApi(fileModel)), nil
}
