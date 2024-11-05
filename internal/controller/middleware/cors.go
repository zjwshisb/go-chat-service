package middleware

import "github.com/gogf/gf/v2/net/ghttp"

func Cors(r *ghttp.Request) {
	corsOptions := r.Response.DefaultCORSOptions()
	corsOptions.AllowDomain = []string{"*"}
	r.Response.CORS(corsOptions)
	r.Middleware.Next()
}
