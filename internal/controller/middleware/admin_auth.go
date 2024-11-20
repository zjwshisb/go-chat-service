package middleware

import (
	"gf-chat/internal/model"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"net/http"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func AdminAuth(r *ghttp.Request) {
	admin, err := service.Admin().Auth(r.GetCtx(), r)
	if err != nil {
		code := gerror.Code(err)
		switch code {
		case gcode.CodeInvalidOperation:
			r.Response.WriteStatus(http.StatusForbidden, g.MapStrStr{
				"message": "forbidden",
			})
			return
		default:
			r.Response.WriteStatus(http.StatusUnauthorized, g.MapStrStr{
				"message": "unauthorized",
			})
		}
		return
	}
	service.AdminCtx().Init(r, &model.AdminCtx{
		Entity: admin,
		Data:   nil,
	})
	r.Middleware.Next()
}
