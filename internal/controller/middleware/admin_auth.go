package middleware

import (
	"gf-chat/internal/dao"
	"gf-chat/internal/model"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"net/http"
)

func getRequestToken(r *ghttp.Request) (token string) {
	token = ""
	bearerToken := r.GetHeader("Authorization")
	if len(bearerToken) > 7 {
		token = bearerToken[7:]
	}
	if token == "" {
		if queryToken := r.Get("token", ""); queryToken.String() != "" {
			token = queryToken.String()
		}
	}
	return token
}

func AdminAuth(r *ghttp.Request) {
	token := getRequestToken(r)
	if token == "" {
		r.Response.WriteStatus(http.StatusUnauthorized)
		return
	}
	uidStr, sessionId, err := service.Jwt().ParseToken(token)
	if err != nil {
		r.Response.WriteStatus(http.StatusUnauthorized)
		return
	}
	uid := gconv.Int(uidStr)
	if uid == 0 {
		r.Response.WriteStatus(http.StatusUnauthorized)
		return
	}
	if sessionId != "" {
		ok := service.Sso().Check(r.Context(), sessionId, uid)
		if !ok {
			r.Response.WriteStatus(http.StatusUnauthorized)
			return
		}
	}

	admin := &entity.CustomerAdmins{}
	dao.CustomerAdmins.Ctx(r.GetCtx()).WherePri(uid).Scan(admin)
	err = service.Admin().IsValid(admin)
	if err != nil {
		r.Response.WriteStatus(http.StatusForbidden, g.MapStrStr{
			"message": err.Error(),
		})
		return
	}
	service.AdminCtx().Init(r, &model.AdminCtx{
		Entity: admin,
		Data:   nil,
	})
	r.Middleware.Next()
}
