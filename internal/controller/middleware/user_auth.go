package middleware

import (
	"gf-chat/internal/model"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/net/ghttp"
	"net/http"
)

func UserAuth(r *ghttp.Request) {
	user, err := service.User().Auth(r.GetCtx(), r)
	if err != nil {
		r.Response.WriteStatus(http.StatusUnauthorized)
		return
	}
	ctx := &model.UserCtx{
		Entity: user,
		Data:   make(map[string]any),
	}
	service.UserCtx().Init(r, ctx)
	r.Middleware.Next()

}
