package backend

import (
	"context"
	"fmt"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

var CImage = &cImage{}

type cImage struct {
}

func (c cImage) Store(ctx context.Context, req *api.ImageReq) (res *api.ImageRes, err error) {
	path := fmt.Sprintf("chat/%d/", service.AdminCtx().GetCustomerId(ctx))
	rPath := req.Path
	if rPath[0:1] == "/" {
		rPath = rPath[1:]
	}
	r, err := service.Qiniu().Save(ctx, req.File, path+rPath)
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, err.Error())
	}
	return &api.ImageRes{
		Url:  r.Url,
		Path: r.Path,
	}, err
}
