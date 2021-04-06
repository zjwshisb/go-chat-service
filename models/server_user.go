package models

import (
	"github.com/gin-gonic/gin"
	"time"
	"ws/db"
	"ws/util"
)
type ServerUserAuthenticate interface {
	Login()
	Logout()
	Auth()
}

type ServerUser struct {
	ID        int64      `json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Username  string     `gorm:"string;size:255" json:"username"`
	Password  string     `gorm:"string;size:255" json:"-"`
	ApiToken string 	`gogm:"string;size:255"  json:"-"`
}

func (user *ServerUser) GetPrimaryKey() int64 {
	return user.ID
}

func (user *ServerUser) Login() (token string) {
	token = util.RandomStr(32)
	db.Db.Model(user).Update("api_token", token)
	return
}

func (user *ServerUser) Logout()  {
	db.Db.Model(user).Update("api_token", "")
}

func (user *ServerUser) Auth(c *gin.Context) {
	db.Db.Where("api_token= ?", util.GetToken(c)).Limit(1).First(user)
}

func (user *ServerUser) FindByName(username string) () {
	db.Db.Where("username= ?", username).Limit(1).First(user)
}
