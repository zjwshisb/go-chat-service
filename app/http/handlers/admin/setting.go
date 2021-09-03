package admin

import (
	"github.com/gin-gonic/gin"
	"ws/app/chat"
	"ws/app/util"
)

type SettingHandler struct {
}

func (handler *SettingHandler) Update(c *gin.Context) {
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
		return
	}
	util.RespSuccess(c, gin.H{})
}

func (handler *SettingHandler) Index(c *gin.Context) {
	var resp = make([]*chat.FieldJson, 0,len(chat.Settings) )
	for _, s := range chat.Settings{
		resp = append(resp, s.ToJson())
	}
	util.RespSuccess(c, resp)
}