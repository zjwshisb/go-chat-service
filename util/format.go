package util

import "github.com/gin-gonic/gin"

func RespSuccess(data interface{}, other ...map[string]interface{}) gin.H {
	resp := gin.H{
		"data": data,
		"success": true,
		"code": 0,
	}
	if len(other) > 0 {
		for _, m:= range other {
			for k, v := range m{
				resp[k] = v
			}
		}
	}

	return resp
}

func RespFail(message string, code int,other ...map[string]interface{}) gin.H {
	resp := gin.H{
		"message": message,
		"success": false,
		"code": code,
	}
	if len(other) > 0 {
		for _, m:= range other {
			for k, v := range m{
				resp[k] = v
			}
		}
	}
	return resp
}
