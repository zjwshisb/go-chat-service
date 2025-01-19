package middleware

import "github.com/gogf/gf/v2/net/ghttp"

func Cors(r *ghttp.Request) {
	r.Response.CORSDefault()
}
