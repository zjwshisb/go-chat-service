package admin

import (
	"github.com/gin-gonic/gin"
	"ws/app/http/requests"
	"ws/app/models"
	"ws/app/util"
	"ws/app/websocket"
)

type UserHandler struct {

}

func (User *UserHandler) Info(c *gin.Context)  {
	admin := requests.GetAdmin(c)
	util.RespSuccess(c, gin.H{
		"username": admin.GetUsername(),
		"id":       admin.GetPrimaryKey(),
	})
}
func (User *UserHandler) Setting(c *gin.Context) {
	u := requests.GetAdmin(c)
	admin := u.(*models.Admin)
	util.RespSuccess(c, admin.GetSetting())
}

func (User *UserHandler) UpdateSetting(c *gin.Context) {
	u := requests.GetAdmin(c)
	admin := u.(*models.Admin)
	form  := requests.AdminChatSettingForm{}
	err := c.ShouldBind(&form)
	if err != nil {
		util.RespValidateFail(c, err.Error())
		return
	}
	setting := admin.GetSetting()
	setting.Background = form.Background
	setting.IsAutoAccept = form.IsAutoAccept
	setting.WelcomeContent = form.WelcomeContent
	setting.OfflineContent = form.OfflineContent
	setting.Name = form.Name
	adminRepo.SaveSetting(setting)
	websocket.AdminManager.PublishUpdateSetting(admin)
	util.RespSuccess(c, gin.H{})
}

func (User *UserHandler) Avatar(c *gin.Context) {
	form := &struct {
		Url string `json:"url"`
	}{}
	err := c.ShouldBind(form)
	if err != nil {
		util.RespError(c, err.Error())
	} else {
		u := requests.GetAdmin(c)
		admin := u.(*models.Admin)
		setting := admin.GetSetting()
		adminRepo.UpdateSetting(setting, "avatar", form.Url)
		websocket.AdminManager.PublishUpdateSetting(admin)
		util.RespSuccess(c, gin.H{})
	}
}

