package backend

import (
	"context"
	"gf-chat/api"
	baseApi "gf-chat/api"
)

var CImage = &cImage{}

type cImage struct {
}

func (c cImage) Store(ctx context.Context, req *api.ImageReq) (res *baseApi.NormalRes[baseApi.ImageRes], err error) {
	//path := fmt.Sprintf("chat/%d/", service.AdminCtx().GetCustomerId(ctx))
	//rPath := req.Path
	//if rPath[0:1] == "/" {
	//	rPath = rPath[1:]
	//}
	return nil, nil
}
