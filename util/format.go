package util

import (
	"github.com/gin-gonic/gin"
)
func RespNotFound(c *gin.Context) {
	c.JSON(404, "404 not found")
}
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
func RespError(c *gin.Context, msg interface{}){
	c.JSON(500, gin.H{
		"message": msg,
	})
}
