package models

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
	"ws/db"
	"ws/util"
)
const (
	User2ServerHashKey = "user-to-server"
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
func (user *User) GetLastServerId() int64 {
	if user.ID == 0 {
		return 0
	}
	ctx := context.Background()
	cmd := db.Redis.HGet(ctx, User2ServerHashKey, string(user.ID))
	if sid, err := cmd.Int64(); err == nil {
		// 判断是否超时
		cmd := db.Redis.ZScore(ctx, fmt.Sprintf(serverChatUserKey, sid), string(user.ID))
		t := int64(cmd.Val())
		if t <= (time.Now().Unix() - 3600 * 24 * 2) {
			return 0
		}
		return sid
	}
	return 0
}