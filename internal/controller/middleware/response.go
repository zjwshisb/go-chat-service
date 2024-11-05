package middleware

import (
	"database/sql"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"net/http"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// SuccessResponse
type SuccessResponse struct {
	Code    int         `json:"code"    dc:"Error code"`
	Data    interface{} `json:"data"    dc:"Result data for certain request according API definition"`
	Success bool        `json:"success" dc:"Is Success"`
}

type FailResponse struct {
	Code    int    `json:"code"    dc:"Error code"`
	Success bool   `json:"success" dc:"Is Success"`
	Message string `json:"message" dc: "错误消息"`
}

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
		if err == sql.ErrNoRows {
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
		r.Response.WriteJson(FailResponse{
			Code:    code.Code(),
			Message: msg,
			Success: false,
		})
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
	r.Response.WriteJson(SuccessResponse{
		Code:    code.Code(),
		Data:    res,
		Success: true,
	})
}
