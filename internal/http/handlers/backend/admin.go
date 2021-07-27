package backend

import (
	"github.com/gin-gonic/gin"
	"ws/internal/auth"
	"ws/internal/databases"
	"ws/internal/file"
	"ws/internal/util"
)

func Me(c *gin.Context) {
	admin := auth.GetAdmin(c)
	util.RespSuccess(c, gin.H{
		"username": admin.Username,
		"id":       admin.ID,
		"avatar":   admin.GetAvatarUrl(),
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
		admin := auth.GetAdmin(c)
		admin.Avatar = fileInfo.Path
		databases.Db.Save(admin)
		util.RespSuccess(c, gin.H{})
	}
}
