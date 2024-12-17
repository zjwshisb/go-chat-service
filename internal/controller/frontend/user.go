package frontend

import (
	"context"
	baseApi "gf-chat/api/v1"
	"gf-chat/api/v1/frontend"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/net/ghttp"
)

var CUser = &cUser{}

type cUser struct {
}

func (c cUser) Login(ctx context.Context, _ *frontend.LoginReq) (res *baseApi.NormalRes[frontend.LoginRes], err error) {
	request := ghttp.RequestFromCtx(ctx)
	_, token, err := service.User().Login(ctx, request)
	if err != nil {
		return
	}
	return baseApi.NewResp(frontend.LoginRes{
		Token: token,
	}), nil

}
