package middleware

import (
	"gf-chat/internal/model"
	"gf-chat/internal/service"
	"net/http"

	"github.com/gogf/gf/v2/net/ghttp"
)

func UserAuth(r *ghttp.Request) {
	token := getRequestToken(r)
	if token == "" {
		r.Response.WriteStatus(http.StatusUnauthorized)
		return
	}
	user := service.User().FindByToken(token)
	if user == nil {
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
