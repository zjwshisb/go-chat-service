package http

import (
	"github.com/gin-gonic/gin"
	"ws/models"
	"ws/util"
)

func Me(c *gin.Context) {
	ui, _ := c.Get("user")
	user := ui.(*models.ServerUser)
	util.RespSuccess(c , gin.H{
		"username": user.Username,
		"id": user.ID,
	})
}
