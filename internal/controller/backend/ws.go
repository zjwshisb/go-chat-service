package backend

import (
	"context"
	baseApi "gf-chat/api"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/net/ghttp"
)

var CWs = &cWs{}

type cWs struct {
}

func (c cWs) Connect(ctx context.Context, req *api.ChatConnectReq) (res *baseApi.NilRes, err error) {
	request := ghttp.RequestFromCtx(ctx)
	conn, err := request.WebSocket()
	if err != nil {
		request.Exit()
		return
	}
	admin := service.AdminCtx().GetAdmin(ctx)
	model := service.Admin().EntityToRelation(admin)
	service.Chat().Register(ctx, model, conn)
	return nil, nil
}
