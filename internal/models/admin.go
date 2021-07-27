package models

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
	"ws/internal/databases"
	"ws/internal/file"
	"ws/internal/util"
)

type AdminJson struct {
	Avatar string `json:"avatar"`
	Username string `json:"username"`
	Online bool `json:"online"`
	Id int64 `json:"id"`
}

type Admin struct {
	ID        int64      `json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Username  string     `gorm:"string;size:255" json:"username"`
	Password  string     `gorm:"string;size:255" json:"-"`
	ApiToken string 	`gorm:"string;size:255"  json:"-"`
	Avatar string 		`gorm:"string;size:512" json:"-"`
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
func (user *Admin) Login() (token string) {
	token = util.RandomStr(32)
	databases.Db.Model(user).Update("api_token", token)
	return
}
func (user *Admin) Logout()  {
	databases.Db.Model(user).Update("api_token", "")
}

func (user *Admin) Auth(c *gin.Context) bool {
	query := databases.Db.Where("api_token= ?", util.GetToken(c)).First(user)
	if query.Error == gorm.ErrRecordNotFound {
		return false
	}
	return true
}
func (user *Admin) FindByName(username string) bool {
	databases.Db.Where("username= ?", username).First(user)
	return user.ID > 0
}
