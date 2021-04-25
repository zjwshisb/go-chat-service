package http

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"path"
	"time"
	"ws/db"
	"ws/models"
	"ws/util"
)

func GetHistoryMessage(c *gin.Context)  {
	ui, _ := c.Get("user")
	user := ui.(*models.User)
	messages := make([]*models.Message, 0)
	db.Db.Preload("ServerUser").Where("user_id = ?", user.ID).
		Where("received_at > ? ", time.Now().Unix() - 2 * 24 * 3600).
		Find(&messages)
	for _, msg := range messages {
		if msg.IsServer {
			msg.Avatar = msg.ServerUser.GetAvatarUrl()
		} else {
			msg.Avatar = ""
		}
	}
	util.RespSuccess(c, messages)
}
// 聊天图片
func Image(c *gin.Context) {
	file, _ := c.FormFile("file")
	ext := path.Ext(file.Filename)
	imageDir := "/chat"
	assetPath := util.StoragePath()
	_, err := os.Stat(assetPath  + imageDir)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(assetPath + imageDir, 0666)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	filename := util.RandomStr(32) + ext
	fullPath := assetPath + imageDir + "/" + filename
	err = c.SaveUploadedFile(file, fullPath)
	if err != nil {
		util.RespError(c, err.Error())
	} else {
		util.RespSuccess(c, gin.H{
			"url": util.Asset( imageDir + "/" + filename),
		})
	}
}