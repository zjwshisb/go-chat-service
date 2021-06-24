package backend

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"ws/internal/databases"
	"ws/internal/file"
	"ws/internal/models"
	"ws/internal/repositories"
	"ws/internal/util"
)

type autoMessageForm struct {
	Name string `form:"name" binding:"required,max=32"`
	Type string `form:"type" binding:"required,autoMessageType"`
	Content string `form:"content" binding:"required,max=512"`
	Title string `form:"title" binding:"max=32"`
	Url string `form:"url" binding:"max=512"`
}

func StoreAutoMessageImage(c *gin.Context) {
	f, _ := c.FormFile("file")
	ff, err := file.Save(f, "auto_message")
	if err != nil {
		util.RespFail(c, err.Error(), 500)
	} else {
		util.RespSuccess(c, gin.H{
			"url": ff.FullUrl,
		})
	}
}
func GetAutoMessages(c *gin.Context)  {
	pagination := repositories.GetAutoMessage([]*repositories.Where{},
	c.GetInt("current"),
	c.GetInt("pageSize"))
	util.RespPagination(c , pagination)
}

func ShowAutoMessage(c *gin.Context) {
	
}

func StoreAutoMessage(c *gin.Context)  {
	form := autoMessageForm{}
	err := c.ShouldBind(&form)
	if err != nil {
		util.RespValidateFail(c, err.Error())
		return
	}
	message := models.AutoMessage{
		Name: form.Name,
		Type: form.Type,
	}
	if message.Type == models.TypeText  || message.Type == models.TypeImage {
		message.Content = form.Content
	}
	if message.Type == models.TypeNavigate {
		content := map[string]string{
			"title": form.Title,
			"url": form.Url,
			"content": form.Content,
		}
		jsonBytes, err := json.Marshal(content)
		if err != nil{
			util.RespError(c, err.Error())
			return
		}
		message.Content = string(jsonBytes)
	}
	databases.Db.Save(&message)
	util.RespSuccess(c, gin.H{})
}
func UpdateAutoMessage(c *gin.Context) {

}
func DeleteAutoMessage(c *gin.Context) {

}
