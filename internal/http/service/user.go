package service

import (
	"github.com/gin-gonic/gin"
	"ws/internal/auth"
	"ws/internal/databases"
	"ws/internal/file"
	"ws/util"
)

func Me(c *gin.Context) {
	user := auth.GetBackendUser(c)
	util.RespSuccess(c, gin.H{
		"username": user.Username,
		"id":       user.ID,
		"avatar":   user.GetAvatarUrl(),
	})
}

// 更新头像
func Avatar(c *gin.Context) {
	f, _ := c.FormFile("file")
	storage := file.Disk("local")
	fileInfo, err := storage.Save(f, "avatar")
	if err != nil {
		util.RespError(c, err.Error())
	} else {
		user := auth.GetBackendUser(c)
		user.Avatar = fileInfo.Path
		databases.Db.Save(user)
		util.RespSuccess(c, gin.H{})
	}
}
