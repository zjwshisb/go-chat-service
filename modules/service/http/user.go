package http

import (
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"path"
	"ws/db"
	"ws/hub"
	"ws/models"
	"ws/util"
)

func Me(c *gin.Context) {
	ui, _ := c.Get("user")
	user := ui.(*models.ServerUser)
	util.RespSuccess(c , gin.H{
		"username": user.Username,
		"id": user.ID,
		"avatar": util.Asset(user.Avatar),
	})
}
// 更新头像
func Avatar(c *gin.Context) {
	file, _ := c.FormFile("file")
	ext := path.Ext(file.Filename)
	avatarDir := "/avatar"
	assetPath := util.StoragePath()
	_, err := os.Stat(assetPath + avatarDir)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(assetPath + avatarDir, 0666)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	temp := os.TempDir() +  "/" + util.RandomStr(32) +  ext
	err = c.SaveUploadedFile(file, temp)
	if err != nil {
		util.RespError(c, err.Error())
	} else {
		imagePath := avatarDir + "/" + util.RandomStr(32) + ext
		file, err := imaging.Open(temp)
		if err != nil {
			util.RespError(c, err.Error())
			return
		}
		err = imaging.Save(imaging.Thumbnail(file, 300,300, imaging.CatmullRom),
			assetPath + imagePath)
		if err != nil {
			util.RespError(c, err.Error())
		} else {
			ui, _ := c.Get("user")
			user := ui.(*models.ServerUser)
			user.Avatar = imagePath
			db.Db.Save(user)
			client, exist := hub.Hub.Server.GetClient(user.ID)
			if exist {
				client.User = user
			}
			util.RespSuccess(c, gin.H{})
		}
		os.Remove(temp)
	}
}