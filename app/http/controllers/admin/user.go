package admin

import (
	"ws/app/http/requests"
	"ws/app/http/responses"
	"ws/app/http/websocket"
	"ws/app/models"
	"ws/app/repositories"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
}

func (User *UserHandler) Info(c *gin.Context) {
	admin := requests.GetAdmin(c)
	responses.RespSuccess(c, gin.H{
		"username": admin.GetUsername(),
		"id":       admin.GetPrimaryKey(),
	})
}
func (User *UserHandler) Setting(c *gin.Context) {
	u := requests.GetAdmin(c)
	admin := u.(*models.Admin)
	responses.RespSuccess(c, admin.GetSetting())
}

func (User *UserHandler) UpdateSetting(c *gin.Context) {
	u := requests.GetAdmin(c)
	admin := u.(*models.Admin)
	form := requests.AdminChatSettingForm{}
	err := c.ShouldBind(&form)
	if err != nil {
		responses.RespValidateFail(c, err.Error())
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
	responses.RespSuccess(c, gin.H{})
}

func (User *UserHandler) Avatar(c *gin.Context) {
	form := &struct {
		Url string `json:"url"`
	}{}
	err := c.ShouldBind(form)
	if err != nil {
		responses.RespError(c, err.Error())
	} else {
		u := requests.GetAdmin(c)
		admin := u.(*models.Admin)
		setting := admin.GetSetting()
		repositories.AdminRepo.UpdateSetting(setting, "avatar", form.Url)
		websocket.AdminManager.PublishUpdateSetting(admin)
		responses.RespSuccess(c, gin.H{})
	}
}
