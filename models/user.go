package models

import (
	"github.com/gin-gonic/gin"
	"time"
	"ws/db"
	"ws/util"
)

type UserAuthenticate interface {
	Delivery()
	Auth()
	Login()
}

type User struct {
	ID        int64      `json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	Username  string     `gorm:"string;size:255" json:"username"`
	Password  string     `gorm:"string;size:255" json:"-"`
	ApiToken  string     `gogm:"string;size:255"  json:"-"`
}

func (user *User) Login() (token string) {
	token = util.RandomStr(32)
	db.Db.Model(user).Update("api_token", token)
	return
}
func (user *User) FindByName(username string) () {
	db.Db.Where("username= ?", username).Limit(1).First(user)
}

func (user *User) Logout() {
	db.Db.Model(user).Update("api_token", "")
}

func (user *User) Auth(c *gin.Context) {
	db.Db.Where("api_token= ?", util.GetToken(c)).Limit(1).First(user)
}
