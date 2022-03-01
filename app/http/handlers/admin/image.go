package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"ws/app/file"
	"ws/app/http/requests"
	"ws/app/util"
)

type ImageHandler struct {
	
}

func (handle *ImageHandler) Store(c *gin.Context)  {
	f, _ := c.FormFile("file")
	admin := requests.GetAdmin(c)
	path := c.Query("path")
	if path == "" {
		util.RespValidateFail(c, "invalid path")
		return
	}
	prefix := fmt.Sprintf("chat/%d/" , admin.GetGroupId())
	ff, err := file.Save(f, prefix + path)
	if err != nil {
		util.RespFail(c, err.Error(), 500)
	} else {
		util.RespSuccess(c, gin.H{
			"url": ff.FullUrl,
		})
	}
}