package admin

import (
	"github.com/gin-gonic/gin"
	"ws/app/auth"
	"ws/app/file"
	"ws/app/http/requests"
	"ws/app/models"
	"ws/app/util"
	"ws/app/websocket"
)

type UserHandler struct {

}

func (User *UserHandler) Info(c *gin.Context)  {
	admin := auth.GetAdmin(c)
	util.RespSuccess(c, gin.H{
		"username": admin.GetUsername(),
		"id":       admin.GetPrimaryKey(),
		"avatar":   admin.GetAvatarUrl(),
	})
}
func (User *UserHandler) Setting(c *gin.Context) {
	admin := auth.GetAdmin(c)
	util.RespSuccess(c, admin.GetSetting())
}
func (User *UserHandler) UpdateSetting(c *gin.Context) {
	admin := auth.GetAdmin(c)
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
	// 如果当前在线，更新信息
	connI , exist := websocket.AdminManager.GetConn(admin)
	if exist {
		admin := connI.GetUser().(*models.Admin)
		admin.RefreshSetting()
	}
	util.RespSuccess(c, gin.H{})
}

func (User *UserHandler) SettingImage(c *gin.Context) {
	f, _ := c.FormFile("file")
	ff, err := file.Save(f, "chat-settings")
	if err != nil {
		util.RespFail(c, err.Error(), 500)
	} else {
		util.RespSuccess(c, gin.H{
			"url": ff.FullUrl,
		})
	}
}

func (User *UserHandler) Avatar(c *gin.Context) {
	f, _ := c.FormFile("file")
	storage := file.Disk("local")
	fileInfo, err := storage.Save(f, "avatar")
	if err != nil {
		util.RespError(c, err.Error())
	} else {
		admin := auth.GetAdmin(c)
		admin.Avatar = fileInfo.Path
		adminRepo.Save(admin)
		// 如果当前在线，更新信息
		conn , exist := websocket.AdminManager.GetConn(admin)
		if exist {
			u := conn.GetUser().(*models.Admin)
			u.RefreshSetting()
		}
		util.RespSuccess(c, gin.H{})
	}
}

