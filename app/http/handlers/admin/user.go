package admin

import (
	"github.com/gin-gonic/gin"
	"ws/app/http/requests"
	"ws/app/models"
	"ws/app/repositories"
	"ws/app/websocket"
)

type UserHandler struct {
}

func (User *UserHandler) Info(c *gin.Context) {
	admin := requests.GetAdmin(c)
	requests.RespSuccess(c, gin.H{
		"username": admin.GetUsername(),
		"id":       admin.GetPrimaryKey(),
	})
}
func (User *UserHandler) Setting(c *gin.Context) {
	u := requests.GetAdmin(c)
	admin := u.(*models.Admin)
	requests.RespSuccess(c, admin.GetSetting())
}

func (User *UserHandler) UpdateSetting(c *gin.Context) {
	u := requests.GetAdmin(c)
	admin := u.(*models.Admin)
	form := requests.AdminChatSettingForm{}
	err := c.ShouldBind(&form)
	if err != nil {
		requests.RespValidateFail(c, err.Error())
		return
	}
	setting := admin.GetSetting()
	setting.Background = form.Background
	setting.IsAutoAccept = form.IsAutoAccept
	setting.WelcomeContent = form.WelcomeContent
	setting.OfflineContent = form.OfflineContent
	setting.Name = form.Name
	repositories.AdminRepo.SaveSetting(setting)
	websocket.AdminManager.PublishUpdateSetting(admin)
	requests.RespSuccess(c, gin.H{})
}

func (User *UserHandler) Avatar(c *gin.Context) {
	form := &struct {
		Url string `json:"url"`
	}{}
	err := c.ShouldBind(form)
	if err != nil {
		requests.RespError(c, err.Error())
	} else {
		u := requests.GetAdmin(c)
		admin := u.(*models.Admin)
		setting := admin.GetSetting()
		repositories.AdminRepo.UpdateSetting(setting, "avatar", form.Url)
		websocket.AdminManager.PublishUpdateSetting(admin)
		requests.RespSuccess(c, gin.H{})
	}
}

