package models

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
	"time"
	"ws/app/databases"
	"ws/app/util"
)




type User struct {
	ID        int64
	CreatedAt *time.Time
	UpdatedAt *time.Time
	Username  string `gorm:"string;size:255" json:"username"`
	Password  string `gorm:"string;size:255" json:"-"`
	ApiToken  string `gogm:"string;size:255"  json:"-"`
	OpenId    string `gorm:"string;size:255"`
}


func (user *User) GetReqId() string {
	key := fmt.Sprintf("user:%d:req-id", user.ID)
	ctx := context.Background()
	cmd := databases.Redis.Incr(ctx, key)
	return "u" + strconv.FormatInt(cmd.Val(), 10)
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
