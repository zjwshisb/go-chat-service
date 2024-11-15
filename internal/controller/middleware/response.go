package middleware

import (
	"database/sql"
	"errors"
	"fmt"
	"gf-chat/api"
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
		}
		// 404错误
		if code == gcode.CodeNotFound {
			r.Response.WriteStatus(http.StatusNotFound, g.MapStrStr{
				"message": msg,
			})
			return
		}
		// 校验错误
		if code == gcode.CodeValidationFailed {
			r.Response.WriteStatus(http.StatusUnprocessableEntity, g.MapStrStr{
				"message": msg,
			})
			return
		}
		// 内部错误
		if code == gcode.CodeNil {
			code = gcode.CodeInternalError
			fmt.Printf("%v", err)
		}
		// 业务错误
		if code == gcode.CodeBusinessValidationFailed {

		}
		r.Response.WriteJson(api.NewFailResp(msg, code.Code()))
		return
	} else if r.Response.Status > 0 && r.Response.Status != http.StatusOK {
		msg = http.StatusText(r.Response.Status)
		switch r.Response.Status {
		case http.StatusNotFound:
			code = gcode.CodeNotFound
		case http.StatusForbidden:
			code = gcode.CodeNotAuthorized
		default:
			code = gcode.CodeUnknown
		}
		r.Response.WriteStatus(r.Response.Status, msg)

		return
	} else {
		code = gcode.CodeOK
	}
	r.Response.WriteJson(res)

}
