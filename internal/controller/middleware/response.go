package middleware

import (
	"database/sql"
	"errors"
	"gf-chat/api/v1"
	"net/http"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// SuccessResponse

// MiddlewareHandlerResponse is the default middleware handling handler response object and its error.
func HandlerResponse(r *ghttp.Request) {
	r.Middleware.Next()

	// There's custom buffer content, it then exits current handler.
	if r.Response.BufferLength() > 0 {
		return
	}

	var (
		msg  string
		err  = r.GetError()
		res  = r.GetHandlerResponse()
		code = gerror.Code(err)
	)
	if err != nil {
		msg = err.Error()
		if errors.Is(err, sql.ErrNoRows) {
			r.Response.WriteStatus(http.StatusNotFound, g.MapStrStr{
				"message": "not found",
			})
			return
		} else if code == gcode.CodeNotFound {
			// 404错误
			r.Response.WriteStatus(http.StatusNotFound, g.MapStrStr{
				"message": msg,
			})
			return
		} else if code == gcode.CodeValidationFailed {
			// 校验错误
			r.Response.WriteStatus(http.StatusUnprocessableEntity, g.MapStrStr{
				"message": msg,
			})
			return
		} else if code == gcode.CodeBusinessValidationFailed {
			// 业务错误
		} else {
			// 非正常错误，记录一下
			g.Log().Error(r.Context(), err)
		}
		r.Response.WriteJson(v1.NewFailResp(msg, code.Code()))
		return
	} else if r.Response.Status > 0 && r.Response.Status != http.StatusOK {
		msg = http.StatusText(r.Response.Status)
		r.Response.WriteStatus(r.Response.Status, msg)

		return
	} else {
		code = gcode.CodeOK
	}
	r.Response.WriteJson(res)

}
