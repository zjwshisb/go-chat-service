package admin

import (
	"github.com/gin-gonic/gin"
	"ws/app/databases"
	"ws/app/http/requests"
	"ws/app/models"
	"ws/app/resource"
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
	admin := requests.GetAdmin(c)
	id := c.Param("id")
	var setting = &models.ChatSetting{}
	databases.Db.Where("group_id = ?" , admin.GetGroupId()).Find(setting, id)
	if setting.Id <= 0 {
		util.RespNotFound(c)
		return
	}
	setting.Value = form.Value
	databases.Db.Save(setting)

	util.RespSuccess(c, gin.H{})
}

func (handler *SettingHandler) Index(c *gin.Context) {
	admin := requests.GetAdmin(c)
	settings := make([]*models.ChatSetting, 0)
	databases.Db.Where("group_id = ?", admin.GetGroupId()).Find(&settings)
	resp := make([]*resource.ChatSetting, len(settings), len(settings))
	for index, setting := range settings {
		resp[index] = setting.ToJson()
	}
	util.RespSuccess(c, resp)
}