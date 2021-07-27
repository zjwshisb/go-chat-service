package admin

import (
	"github.com/gin-gonic/gin"
	"ws/internal/chat"
	"ws/internal/json"
	"ws/internal/util"
)


func UpdateSetting(c *gin.Context) {
	var form = struct {
		Value string `json:"value" form:"value" binding:"required"`
	}{}
	err := c.ShouldBind(&form)
	if err != nil {
		util.RespValidateFail(c, err.Error())
		return
	}
	name := c.Param("name")
	setting, exist := chat.Settings[name]
	if !exist {
		util.RespValidateFail(c, err.Error())
		return
	}
	err = setting.SetValue(form.Value)
	if err !=nil {
		util.RespValidateFail(c , err.Error())
	}
	util.RespSuccess(c, gin.H{})
}

func GetSettings(c *gin.Context) {
	var resp = make([]json.Setting, 0,len(chat.Settings) )
	for _, s := range chat.Settings{
		resp = append(resp, json.Setting{
			Name:    s.Name,
			Title:   s.Title,
			Value:   s.GetValue(),
			Options: s.Options,
		})
	}
	util.RespSuccess(c, resp)
}