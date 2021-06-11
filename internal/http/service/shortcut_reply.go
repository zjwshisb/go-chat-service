package service

import (
	"github.com/gin-gonic/gin"
	"ws/internal/databases"
	"ws/internal/models"
	"ws/util"
)

type ShortcutForm struct {
	Content string `form:"content" binding:"required"`
}


func UpdateShortcutReply(c *gin.Context){
	user := getUser(c)
	util.RespSuccess(c, user.ShortcutReplies)
}

func StoreShortcutReply(c *gin.Context)  {
	form := &ShortcutForm{}
	user := getUser(c)
	err := c.ShouldBind(form)
	if err != nil {
		util.RespValidateFail(c, "表单验证失败")
		return
	}
	reply := &models.ShortcutReply{
		Content: form.Content,
		UserId: user.ID,
	}
	databases.Db.Save(reply)
	util.RespSuccess(c, gin.H{})
}
func DeleteShortcutReply(c *gin.Context)  {
	user := getUser(c)
	databases.Db.
		Where("id = ?" ,c.Param("id")).
		Where("user_id= ?", user.ID).
		Delete(&models.ShortcutReply{})
	util.RespSuccess(c, gin.H{})
}

func GetShortcutReply(c *gin.Context) {
	user := getUser(c)
	replies := make([]*models.ShortcutReply, 0)
	databases.Db.Model(user).Association("ShortcutReplies").Find(&replies)
	util.RespSuccess(c, replies)
}