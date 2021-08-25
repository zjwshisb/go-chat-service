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
	Username      string `json:"username"`
	Online        bool   `json:"online"`
	Id            int64  `json:"id"`
	AcceptedCount int64  `json:"accepted_count"`
}

type Admin struct {
	ID        int64
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
	Username  string            `gorm:"string;size:255" `
	Password  string            `gorm:"string;size:255" `
	ApiToken  string            `gorm:"string;size:255" `
	Avatar    string            `gorm:"string;size:512"`
	Setting   *AdminChatSetting `json:"message" gorm:"foreignKey:admin_id"`
}

func (user *Admin) AccessTo(uid int64) bool {
	return true
}

func (user *Admin) GetPrimaryKey() int64 {
	return user.ID
}
func (user *Admin) GetAvatarUrl() string {
	if user.Avatar != "" {
		return file.Disk("local").Url(user.Avatar)
	}
	return ""
}
func (user *Admin) GetUsername() string {
	return user.Username
}
func (user *Admin) GetChatName() string {
	if user.Setting == nil {
		setting := &AdminChatSetting{}
		databases.Db.Model(user).Association("Setting").Find(setting)
		user.Setting = setting
	}
	if user.Setting != nil {
		if user.Setting.Name != "" {
			return user.Setting.Name
		}
	}
	return user.Username
}
func (user *Admin) Login() (token string) {
	token = util.RandomStr(32)
	databases.Db.Model(user).Update("api_token", token)
	return
}
func (user *Admin) Logout() {
	databases.Db.Model(user).Update("api_token", "")
}

func (user *Admin) Auth(c *gin.Context) bool {
	token := util.GetToken(c)
	if token == "" {
		return false
	}
	query := databases.Db.Where("api_token= ?", token).First(user)
	if query.Error == gorm.ErrRecordNotFound {
		return false
	}
	return true
}
func (user *Admin) FindByName(username string) bool {
	databases.Db.Where("username= ?", username).First(user)
	return user.ID > 0
}
