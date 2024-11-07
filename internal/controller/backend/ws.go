package backend

import (
	"context"
	baseApi "gf-chat/api"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/service"
	"net/http"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gorilla/websocket"
)

var CWs = &cWs{}
var (
	update = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type cWs struct {
}

func (c cWs) Connect(ctx context.Context, req *api.ChatConnectReq) (res *baseApi.NilRes, err error) {
	request := ghttp.RequestFromCtx(ctx)
	conn, err := update.Upgrade(request.Response.Writer, request.Request, nil)
	if err != nil {
		request.Exit()
		return
	}
	admin := service.AdminCtx().GetAdmin(ctx)
	service.Chat().Register(ctx, admin, conn)
	return nil, nil
}
