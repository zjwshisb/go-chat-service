package user

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"ws/internal/databases"
	"ws/internal/file"
	"ws/internal/models"
	"ws/internal/resources"
	"ws/util"
)

// 消息记录
func GetHistoryMessage(c *gin.Context) {
	ui, _ := c.Get("user")
	user := ui.(*models.User)
	messages := make([]*models.Message, 0)
	query := databases.Db.Preload("ServerUser").Where("user_id = ?", user.ID)
	id, exist := c.GetQuery("id")
	if exist {
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err == nil {
			query.Where("id < ?", idInt)
		}
	}
	query.Order("id desc").Limit(20).Find(&messages)
	messagesResources := make([]*resources.Message, 0)
	for _, msg := range messages {
		messagesResources = append(messagesResources, resources.NewMessage(*msg))
	}
	util.RespSuccess(c, messagesResources)
}

// 聊天图片
func Image(c *gin.Context) {
	f, _ := c.FormFile("file")
	ff, err := file.Save(f, "chat")
	if err != nil {
		util.RespFail(c, err.Error(), 500)
	} else {
		util.RespSuccess(c, gin.H{
			"url": ff.FullUrl,
		})
	}
}
