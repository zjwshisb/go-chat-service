package admin

import (
	"fmt"
	"ws/app/file"
	"ws/app/http/requests"
	"ws/app/http/responses"

	"github.com/gin-gonic/gin"
)

type ImageHandler struct {
}

func (handle *ImageHandler) Store(c *gin.Context) {
	f, _ := c.FormFile("file")
	admin := requests.GetAdmin(c)
	path := c.Query("path")
	if path == "" {
		responses.RespValidateFail(c, "invalid path")
		return
	}
	prefix := fmt.Sprintf("chat/%d/", admin.GetGroupId())
	ff, err := file.Save(f, prefix+path)
	if err != nil {
		responses.RespFail(c, err.Error(), 500)
	} else {
		responses.RespSuccess(c, gin.H{
			"url": ff.FullUrl,
		})
	}
}
