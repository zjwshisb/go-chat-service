package http

import (
	"github.com/gin-gonic/gin"
	"time"
	"ws/db"
	"ws/models"
	"ws/util"
)

func GetHistoryMessage(c *gin.Context)  {
	ui, _ := c.Get("user")
	user := ui.(*models.User)
	messages := make([]*models.Message, 0)
	db.Db.Where("user_id = ?", user.ID).
		Where("received_at > ? ", time.Now().Unix() - 2 * 24 * 3600).
		Find(&messages)
	util.RespSuccess(c, messages)
}
