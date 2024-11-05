package frontend

import (
	"context"
	"fmt"
	baseApi "gf-chat/api"
	"gf-chat/api/v1/frontend"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
)

var CWs = &cWs{}

type cWs struct {
}

func (c cWs) GetReqId(ctx context.Context, req *frontend.ChatReqIdReq) (res *frontend.ChatReqIdRes, err error) {
	return &frontend.ChatReqIdRes{ReqId: service.ChatMessage().GenReqId()}, nil
}

func (c cWs) Image(ctx context.Context, req *frontend.ChatImageReq) (res *frontend.ChatImageRes, err error) {
	path := fmt.Sprintf("chat/%d/user", service.UserCtx().GetCustomerId(ctx))
	r, err := service.Qiniu().Save(ctx, req.File, path)
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, err.Error())
	}
	return &frontend.ChatImageRes{
		Url:  r.Url,
		Path: r.Path,
	}, nil
}

func (c cWs) Connect(ctx context.Context, req *frontend.ChatConnectReq) (res *baseApi.NilRes, err error) {
	request := ghttp.RequestFromCtx(ctx)
	conn, err := request.WebSocket()
	if err != nil {
		request.Exit()
		return
	}
	user := service.UserCtx().GetUser(ctx)
	service.Chat().Register(ctx, user, conn)
	return nil, nil
}
