package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"ws/app/file"
	"ws/app/http/requests"
)

type ImageHandler struct {
}

func (handle *ImageHandler) Store(c *gin.Context) {
	f, _ := c.FormFile("file")
	admin := requests.GetAdmin(c)
	path := c.Query("path")
	if path == "" {
		requests.RespValidateFail(c, "invalid path")
		return
	}
	prefix := fmt.Sprintf("chat/%d/", admin.GetGroupId())
	ff, err := file.Save(f, prefix+path)
	if err != nil {
		requests.RespFail(c, err.Error(), 500)
	} else {
		requests.RespSuccess(c, gin.H{
			"url": ff.FullUrl,
		})
	}
}