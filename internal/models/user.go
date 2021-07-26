package models

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
	"ws/internal/databases"
	"ws/internal/util"
)


type UserJson struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	LastChatTime int64  `json:"last_chat_time"`
	Disabled bool       `json:"disabled"`
	Online bool         `json:"online"`
	Messages []*MessageJson `json:"messages"`
	Unread int          `json:"unread"`
}


type User struct {
	ID        int64      `json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	Username  string     `gorm:"string;size:255" json:"username"`
	Password  string     `gorm:"string;size:255" json:"-"`
	ApiToken  string     `gogm:"string;size:255"  json:"-"`
	OpenId string `gorm:"string;size:255"`
}


func (user *User) GetUsername() string  {
	return user.Username
}
func (user *User) GetAvatarUrl() string {
	return ""
}
func (user *User) GetPrimaryKey() int64 {
	return user.ID
}
func (user *User) GetMpOpenId() string {
	return user.OpenId
}

func (user *User) Auth(c *gin.Context) bool {
	query := databases.Db.Where("api_token= ?", util.GetToken(c)).Limit(1).First(user)
	if query.Error == gorm.ErrRecordNotFound {
		return false
	}
	return true
}

func (user *User) Login() (token string) {
	token = util.RandomStr(32)
	databases.Db.Model(user).Update("api_token", token)
	return
}
func (user *User) FindByName(username string) {
	databases.Db.Where("username= ?", username).Limit(1).First(user)
}

