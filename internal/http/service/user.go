package service

import (
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"os"
	"path"
	"ws/internal/databases"
	"ws/internal/image"
	"ws/internal/models"
	"ws/util"
)

func Me(c *gin.Context) {
	ui, _ := c.Get("user")
	user := ui.(*models.ServiceUser)
	util.RespSuccess(c, gin.H{
		"username": user.Username,
		"id":       user.ID,
		"avatar":   image.Url(user.Avatar),
	})
}

// 更新头像
func Avatar(c *gin.Context) {
	file, _ := c.FormFile("file")
	ext := path.Ext(file.Filename)
	var err error
	temp := os.TempDir() + "/" + util.RandomStr(32) + ext
	err = c.SaveUploadedFile(file, temp)
	if err != nil {
		util.RespError(c, err.Error())
	} else {
		imagePath := image.AvatarDIR + "/" + util.RandomStr(32) + ext
		file, err := imaging.Open(temp)
		if err != nil {
			util.RespError(c, err.Error())
			return
		}
		err = imaging.Save(imaging.Thumbnail(file, 300, 300, imaging.CatmullRom),
			image.BasePath+imagePath)
		if err != nil {
			util.RespError(c, err.Error())
		} else {
			ui, _ := c.Get("user")
			user := ui.(*models.ServiceUser)
			user.Avatar = imagePath
			databases.Db.Save(user)
			util.RespSuccess(c, gin.H{})
		}
		os.Remove(temp)
	}
}
