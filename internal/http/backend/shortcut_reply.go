package backend

import (
	"github.com/gin-gonic/gin"
	"ws/internal/auth"
	"ws/internal/repositories"
	"ws/util"
)

type ShortcutForm struct {
	Content string `form:"content" binding:"required,max=255"`
}

func UpdateShortcutReply(c *gin.Context){
	user := auth.GetBackendUser(c)
	util.RespSuccess(c, user.ShortcutReplies)
}

func StoreShortcutReply(c *gin.Context)  {
	form := &ShortcutForm{}
	user := auth.GetBackendUser(c)
	err := c.ShouldBind(form)
	if err != nil {
		util.RespValidateFail(c, "表单验证失败")
		return
	}
	repositories.StoreShortcutReply(map[string]interface{}{
		"Content": form.Content,
		"UserId": user.GetPrimaryKey(),
	})
	util.RespSuccess(c, gin.H{})
}

func DeleteShortcutReply(c *gin.Context)  {
	user := auth.GetBackendUser(c)
	rowsAffected := repositories.DeleteShortcutReply([]repositories.Where{
		{
			Filed: "user_id = ?",
			Value: user.GetPrimaryKey(),
		},
		{
			Filed: "id = ?",
			Value: c.Param("id"),
		},
	})
	if rowsAffected == 0 {
		util.RespNotFound(c)
	} else {
		util.RespSuccess(c, gin.H{})
	}
}

func GetShortcutReply(c *gin.Context) {
	user := auth.GetBackendUser(c)
	replies := repositories.GetShortcutReply([]repositories.Where{
		{
			Filed: "user_id = ?",
			Value: user.GetPrimaryKey(),
		},
	})
	util.RespSuccess(c, replies)
}