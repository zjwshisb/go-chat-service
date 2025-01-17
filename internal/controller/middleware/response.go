package middleware

import (
	"database/sql"
	"errors"
	"gf-chat/api"
	"net/http"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// SuccessResponse

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
		if errors.Is(err, sql.ErrNoRows) || code == gcode.CodeNotFound {
			r.Response.WriteStatus(http.StatusNotFound, g.MapStrStr{
				"message": "not found",
			})
			return
		} else if code == gcode.CodeValidationFailed {
			// 校验错误
			r.Response.WriteStatus(http.StatusUnprocessableEntity, g.MapStrStr{
				"message": msg,
			})
			return
		} else if code == gcode.CodeBusinessValidationFailed {
			r.Response.WriteJson(api.NewFailResp(msg, code.Code()))
			// 业务错误
		} else {
			g.Log().Errorf(r.Context(), "%+v", err)
		}
		r.Response.WriteStatus(http.StatusInternalServerError, g.MapStrStr{
			"message": "internal server error",
		})
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
