package models

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
	"ws/app/databases"
	"ws/app/file"
	"ws/app/util"
)

type AdminJson struct {
	Avatar        string `json:"avatar"`
	adminname      string `json:"adminname"`
	Online        bool   `json:"online"`
	Id            int64  `json:"id"`
	AcceptedCount int64  `json:"accepted_count"`
}

type Admin struct {
	ID        int64
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
	username  string            `gorm:"string;size:255" `
	Password  string            `gorm:"string;size:255" `
	ApiToken  string            `gorm:"string;size:255" `
	Avatar    string            `gorm:"string;size:512"`
	Setting   *AdminChatSetting `json:"message" gorm:"foreignKey:admin_id"`
	IsSuper bool `gorm:"is_super"`
}
// 是否有admin的权限
func (admin *Admin) AccessTo(uid int64) bool {
	return true
}
func (admin *Admin)  GetIsSuper() bool {
	return admin.IsSuper
}

func (admin *Admin) GetPrimaryKey() int64 {
	return admin.ID
}
func (admin *Admin) GetAvatarUrl() string {
	if admin.Avatar != "" {
		return file.Disk("local").Url(admin.Avatar)
	}
	return util.SystemAvatar()
}
func (admin *Admin) GetUsername() string {
	return admin.username
}
// 客服名称
func (admin *Admin) GetChatName() string {
	if admin.Setting == nil {
		setting := &AdminChatSetting{}
		databases.Db.Model(admin).Association("Setting").Find(setting)
		admin.Setting = setting
	}
	if admin.Setting != nil {
		if admin.Setting.Name != "" {
			return admin.Setting.Name
		}
	}
	return admin.GetUsername()
}
func (admin *Admin) Login() (token string) {
	token = util.RandomStr(32)
	databases.Db.Model(admin).Update("api_token", token)
	return
}
func (admin *Admin) Logout() {
	databases.Db.Model(admin).Update("api_token", "")
}

func (admin *Admin) Auth(c *gin.Context) bool {
	token := util.GetToken(c)
	if token == "" {
		return false
	}
	query := databases.Db.Where("api_token= ?", token).First(admin)
	if query.Error == gorm.ErrRecordNotFound {
		return false
	}
	return true
}
func (admin *Admin) FindByName(username string) bool {
	databases.Db.Where("username= ?", username).First(admin)
	return admin.ID > 0
}
