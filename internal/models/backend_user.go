package models

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
	"ws/internal/databases"
	"ws/internal/file"
	"ws/internal/util"
)

type BackendUser struct {
	ID        int64      `json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Username  string     `gorm:"string;size:255" json:"username"`
	Password  string     `gorm:"string;size:255" json:"-"`
	ApiToken string 	`gorm:"string;size:255"  json:"-"`
	Avatar string 		`gorm:"string;size:512" json:"-"`
}

func (user *BackendUser) GetPrimaryKey() int64 {
	return user.ID
}
func (user *BackendUser) GetAvatarUrl() string {
	if user.Avatar != "" {
		return file.Disk("local").Url(user.Avatar)
	}
	return ""
}
func (user *BackendUser) Login() (token string) {
	token = util.RandomStr(32)
	databases.Db.Model(user).Update("api_token", token)
	return
}
func (user *BackendUser) Logout()  {
	databases.Db.Model(user).Update("api_token", "")
}

func (user *BackendUser) Auth(c *gin.Context) bool {
	query := databases.Db.Where("api_token= ?", util.GetToken(c)).First(user)
	if query.Error == gorm.ErrRecordNotFound {
		return false
	}
	return true
}
func (user *BackendUser) FindByName(username string) bool {
	databases.Db.Where("username= ?", username).First(user)
	return user.ID > 0
}
