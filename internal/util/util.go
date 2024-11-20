package util

import "github.com/gogf/gf/v2/net/ghttp"

func GetRequestToken(r *ghttp.Request) (token string) {
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
