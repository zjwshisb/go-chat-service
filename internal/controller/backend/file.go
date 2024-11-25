package backend

import (
	"context"
	baseApi "gf-chat/api"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/consts"
	"gf-chat/internal/library/storage"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

var CImage = &cFile{}

type cFile struct {
}

func (c cFile) Index(ctx context.Context, req *api.FileListReq) (res *baseApi.ListRes[*api.File], err error) {
	where := do.CustomerChatFiles{
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
		ParentId:   req.DirId,
	}
	files, err := service.File().All(ctx, where, nil, "id desc")
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

func (c cFile) StoreDir(ctx context.Context, req *api.FileDirStoreReq) (res *baseApi.NilRes, err error) {
	var parent *model.CustomerChatFile
	if req.Pid != 0 {
		parent, err = service.File().First(ctx, do.CustomerChatFiles{
			Id:         req.Pid,
			CustomerId: service.AdminCtx().GetCustomerId(ctx),
			Type:       consts.FileTypeDir,
		})
		if err != nil {
			return nil, err
		}
		if parent == nil {
			return nil, gerror.NewCode(gcode.CodeValidationFailed, "dir not exists")
		}
	}
	dirExists, err := service.File().Exists(ctx, do.CustomerChatFiles{
		Name:       req.Name,
		ParentId:   req.Pid,
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
	})
	if err != nil {
		return
	}
	if dirExists {
		return nil, gerror.NewCode(gcode.CodeValidationFailed, "已存在同名文件夹")
	}
	path := req.Name
	if parent != nil {
		path = parent.Path + "/" + path
	}
	newDir := model.CustomerChatFile{
		CustomerChatFiles: entity.CustomerChatFiles{
			Name:       req.Name,
			ParentId:   req.Pid,
			Path:       path,
			Disk:       "",
			Type:       consts.FileTypeDir,
			FromId:     service.AdminCtx().GetId(ctx),
			FromModel:  "admin",
			CustomerId: service.AdminCtx().GetCustomerId(ctx),
		},
	}
	_, err = service.File().Save(ctx, &newDir)
	if err != nil {
		return
	}
	return
}
