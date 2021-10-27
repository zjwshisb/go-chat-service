package admin

import (
	"github.com/gin-gonic/gin"
	"ws/app/chat"
	"ws/app/json"
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
	setting, exist := chat.SettingService.Values[name]
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
	var resp = make([]*json.SettingField, 0,len(chat.SettingService.Values) )
	for _, s := range chat.SettingService.Values{
		resp = append(resp, s.ToJson())
	}
	util.RespSuccess(c, resp)
}