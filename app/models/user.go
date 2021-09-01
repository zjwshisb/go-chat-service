package models

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
	"ws/app/databases"
	"ws/app/util"
)

type UserJson struct {
	ID           int64          `json:"id"`
	Username     string         `json:"username"`
	LastChatTime int64          `json:"last_chat_time"`
	Disabled     bool           `json:"disabled"`
	Online       bool           `json:"online"`
	Messages     []*MessageJson `json:"messages"`
	Unread       int            `json:"unread"`
	Avatar string `json:"avatar"`
}

type WaitingUserJson struct {
	Username     string `json:"username"`
	Avatar       string `json:"avatar"`
	Id           int64  `json:"id"`
	LastMessage  string `json:"last_message"`
	LastTime     int64  `json:"last_time"`
	LastType     string `json:"last_type"`
	MessageCount int    `json:"message_count"`
	Description  string `json:"description"`
}

type User struct {
	ID        int64
	CreatedAt *time.Time
	UpdatedAt *time.Time
	Username  string `gorm:"string;size:255" json:"username"`
	Password  string `gorm:"string;size:255" json:"-"`
	ApiToken  string `gogm:"string;size:255"  json:"-"`
	OpenId    string `gorm:"string;size:255"`
}

func (user *User) GetUsername() string {
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
	token := util.GetToken(c)
	if token == "" {
		return false
	}
	query := databases.Db.Where("api_token= ?", token).Limit(1).First(user)
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
