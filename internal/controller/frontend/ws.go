package frontend

import (
	"context"
	baseApi "gf-chat/api/v1"
	"gf-chat/api/v1/frontend"
	"gf-chat/internal/service"
	"net/http"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gorilla/websocket"
)

var (
	update = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

var CWs = &cWs{}

type cWs struct {
}

func (c cWs) Image(ctx context.Context, req *frontend.ChatImageReq) (res *frontend.ChatImageRes, err error) {
	//path := fmt.Sprintf("chat/%d/user", service.UserCtx().GetCustomerId(ctx))
	//r, err := service.Qiniu().Save(ctx, req.File, path)
	//if err != nil {
	//	return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, err.Error())
	//}
	//return &frontend.ChatImageRes{
	//	Url:  r.Url,
	//	Path: r.Path,
	//}, nil
	return nil, nil
}

func (c cWs) Connect(ctx context.Context, _ *frontend.ChatConnectReq) (res *baseApi.NilRes, err error) {
	request := ghttp.RequestFromCtx(ctx)
	conn, err := update.Upgrade(request.Response.Writer, request.Request, nil)
	if err != nil {
		request.Exit()
		return
	}
	user := service.UserCtx().GetUser(ctx)
	err = service.Chat().Register(user, conn, service.Platform().GetPlatform(ctx))
	if err != nil {
		return
	}
	res = baseApi.NewNilResp()
	return
}
