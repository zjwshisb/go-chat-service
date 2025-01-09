package frontend

import (
	"context"
	baseApi "gf-chat/api"
	"gf-chat/api/frontend/v1"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/net/ghttp"
)

var CUser = &cUser{}

type cUser struct {
}

func (c cUser) Login(ctx context.Context, _ *v1.LoginReq) (res *baseApi.NormalRes[v1.LoginRes], err error) {
	request := ghttp.RequestFromCtx(ctx)
	_, token, err := service.User().Login(ctx, request)
	if err != nil {
		return
	}
	return baseApi.NewResp(v1.LoginRes{
		Token: token,
	}), nil

}
