package middleware

import (
	"gf-chat/internal/model"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/net/ghttp"
	"net/http"
)

func UserAuth(r *ghttp.Request) {
	token := getRequestToken(r)
	if token == "" {
		r.Response.WriteStatus(http.StatusUnauthorized)
		return
	}
	user, userApp := service.User().FindByToken(token)
	if user == nil {
		r.Response.WriteStatus(http.StatusUnauthorized)
		return
	}
	ctx := &model.UserCtx{
		Entity:  user,
		UserApp: userApp,
		Data:    make(map[string]any),
	}
	service.UserCtx().Init(r, ctx)
	r.Middleware.Next()

}
