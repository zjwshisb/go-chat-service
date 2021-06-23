package backend

import (
	"github.com/gin-gonic/gin"
	"ws/internal/databases"
	"ws/internal/file"
	"ws/internal/models"
	"ws/util"
)
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
	var messages []models.AutoMessage
	databases.Db.Find(messages)
	util.RespSuccess(c , messages)
}
func ShowAutoMessage(c *gin.Context) {

}
func StoreAutoMessage(c *gin.Context)  {
	
}
func UpdateAutoMessage(c *gin.Context) {

}
func DeleteAutoMessage(c *gin.Context) {

}
