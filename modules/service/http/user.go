package http

import (
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"os"
	"path"
	"ws/db"
	"ws/models"
	"ws/util"
)

func Me(c *gin.Context) {
	ui, _ := c.Get("user")
	user := ui.(*models.ServerUser)
	util.RespSuccess(c , gin.H{
		"username": user.Username,
		"id": user.ID,
		"avatar": util.Asset(c, user.Avatar),
	})
}
// 上传头像
func Avatar(c *gin.Context) {
	file, _ := c.FormFile("file")
	ext := path.Ext(file.Filename)
	pwd, _ := os.Getwd()
	avatarDir := "/storage/avatar"
	_, err := os.Stat(pwd + avatarDir)
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(pwd + avatarDir, 0666)
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
			pwd + imagePath)
		if err != nil {
			util.RespError(c, err.Error())
		} else {
			ui, _ := c.Get("user")
			user := ui.(*models.ServerUser)
			user.Avatar = imagePath
			db.Db.Save(user)
			util.RespSuccess(c, gin.H{})
		}
		os.Remove(temp)
	}
}