package backend

import (
	"context"
	baseApi "gf-chat/api/v1"
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
	"github.com/gogf/gf/v2/frame/g"
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
	// 目录排到最前面
	slice.SortBy(apiFiles, func(item1, item2 *api.File) bool {
		if item1.Type == consts.FileTypeDir && item2.Type != consts.FileTypeDir {
			return true
		}
		if item1.Type != consts.FileTypeDir && item2.Type == consts.FileTypeDir {
			return false
		}
		return item1.Id > item2.Id
	})

	return baseApi.NewListResp(apiFiles, count), nil
}

func (c cFile) Store(ctx context.Context, req *api.FileStoreReq) (res *baseApi.NormalRes[*api.File], err error) {
	var parent *model.CustomerChatFile
	request := g.RequestFromCtx(ctx)
	dirVal := request.GetCtxVar("file-dir")
	if p, ok := dirVal.Interface().(*model.CustomerChatFile); ok {
		parent = p
	}
	file, err := req.File.Open()
	if err != nil {
		return
	}
	defer func() {
		_ = file.Close()
	}()
	fileType, _ := storage.FileType(file)
	relativePath := ""
	if parent != nil {
		relativePath = parent.Path
	}
	fileModel, err := storage.Disk().SaveUpload(ctx, req.File, relativePath)
	if err != nil {
		return nil, err
	}
	fileModel.Type = fileType
	fileModel.ParentId = req.Pid
	fileModel.CustomerId = service.AdminCtx().GetCustomerId(ctx)
	fileModel.FromId = service.AdminCtx().GetId(ctx)
	fileModel.FromModel = "admin"
	_, err = service.File().Insert(ctx, fileModel)
	if err != nil {
		return
	}
	return baseApi.NewResp(service.File().ToApi(fileModel)), nil
}
func (c cFile) Update(ctx context.Context, req *api.FileUpdateReq) (res *baseApi.NormalRes[*api.File], err error) {
	file, err := service.File().First(ctx, do.CustomerChatFiles{
		Id:         g.RequestFromCtx(ctx).GetRouter("id").Uint(),
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
	})
	if err != nil {
		return
	}
	existsWhere := g.Map{
		"name": req.Name,
	}
	if file.Type == consts.FileTypeDir {
		existsWhere["type"] = consts.FileTypeDir
	} else {
		existsWhere["type !="] = consts.FileTypeDir
	}
	exist, err := service.File().Exists(ctx, existsWhere)
	if err != nil {
		return
	}
	if exist {
		name := "文件"
		if file.Type == consts.FileTypeDir {
			name = "目录"
		}
		err = gerror.New("存在同名的" + name)
		return
	}
	file.Name = req.Name
	_, err = service.File().Save(ctx, file)
	if err != nil {
		return
	}
	return baseApi.NewResp(service.File().ToApi(file)), nil
}

func (c cFile) StoreDir(ctx context.Context, req *api.FileDirStoreReq) (res *baseApi.NormalRes[*api.File], err error) {
	request := g.RequestFromCtx(ctx)
	var parent *model.CustomerChatFile
	dirVal := request.GetCtxVar("file-dir")
	if p, ok := dirVal.Interface().(*model.CustomerChatFile); ok {
		parent = p
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
		err = gerror.NewCode(gcode.CodeValidationFailed, "已存在同名目录")
		return
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
	return baseApi.NewResp(service.File().ToApi(&newDir)), nil
}
func (c cFile) Delete(ctx context.Context, _ *api.FileDeleteReq) (res *baseApi.NilRes, err error) {
	file, err := service.File().First(ctx, do.CustomerChatFiles{
		Id:         g.RequestFromCtx(ctx).GetRouter("id").Uint(),
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
	})
	if err != nil {
		return
	}
	if file.Type == consts.FileTypeDir {
		childrenExists := false
		childrenExists, err = service.File().Exists(ctx, do.CustomerChatFiles{
			ParentId: file.Id,
		})
		if err != nil {
			return
		}
		if childrenExists {
			err = gerror.New("该目录下存在其他文件，无法删除")
			return
		}
	}
	err = service.File().Delete(ctx, file.Id)

	return baseApi.NewNilResp(), err
}
