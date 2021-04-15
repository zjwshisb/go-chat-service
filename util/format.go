package util

import (
	"github.com/gin-gonic/gin"
)

func RespSuccess(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"data": data,
		"success": true,
		"code": 0,
	})
}
func RespValidateFail(c *gin.Context, msg interface{}) {
	c.JSON(422, gin.H{
		"message": msg,
		"success": false,
	})
}
func RespFail(c *gin.Context, msg interface{} , code int) {
	c.JSON(200, gin.H{
		"success": false,
		"code":code,
		"message": msg,
	})
}

