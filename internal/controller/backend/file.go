package backend

import (
	"context"
	"gf-chat/api"
	baseApi "gf-chat/api"
	"gf-chat/internal/library/storage"
	"github.com/gogf/gf/v2/frame/g"
)

var CImage = &cFile{}

type cFile struct {
}

func (c cFile) Store(ctx context.Context, req *api.FileReq) (res *baseApi.NormalRes[baseApi.FileRes], err error) {
	file := req.File
	name, err := storage.Disk().SaveUpload(ctx, file, "test")
	if err != nil {
		return nil, err
	}
	g.Dump(name)
	//path := fmt.Sprintf("chat/%d/", service.AdminCtx().GetCustomerId(ctx))
	//rPath := req.Path
	//if rPath[0:1] == "/" {
	//	rPath = rPath[1:]
	//}
	return nil, nil
}
