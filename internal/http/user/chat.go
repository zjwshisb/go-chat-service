package user

import (
	"github.com/gin-gonic/gin"
	"path"
	"strconv"
	"ws/core/image"
	"ws/db"
	"ws/internal/models"
	resources2 "ws/internal/resources"
	"ws/util"
)
// 消息记录
func GetHistoryMessage(c *gin.Context)  {
	ui, _ := c.Get("user")
	user := ui.(*models.User)
	messages := make([]*models.Message, 0)
	query := db.Db.Preload("ServerUser").Where("user_id = ?", user.ID)
	id, exist := c.GetQuery("id")
	if exist {
		idInt, err  := strconv.ParseInt(id, 10, 64)
		if err == nil {
			query.Where("id < ?", idInt)
		}
	}
	query.Order("id desc").Limit(20).Find(&messages)
	messagesResources := make([]*resources2.Message,0)
	for _, msg := range messages {
		messagesResources = append(messagesResources, resources2.NewMessage(*msg))
	}
	util.RespSuccess(c, messagesResources)
}
// 聊天图片
func Image(c *gin.Context) {
	file, _ := c.FormFile("file")
	ext := path.Ext(file.Filename)
	filename := util.RandomStr(32) + ext
	fullPath := image.BasePath + image.ChatDir + "/" + filename
	err := c.SaveUploadedFile(file, fullPath)
	if err != nil {
		util.RespError(c, err.Error())
	} else {
		util.RespSuccess(c, gin.H{
			"url": image.Url(image.ChatDir + "/" + filename),
		})
	}
}