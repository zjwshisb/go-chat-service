package models

import (
	"context"
	"fmt"
	"github.com/duke-git/lancet/v2/random"
	"strconv"
	"time"
	"ws/app/contract"
	"ws/app/databases"
)

type User struct {
	ID        int64
	CreatedAt *time.Time
	UpdatedAt *time.Time
	Username  string `gorm:"string;size:255" json:"username"`
	Password  string `gorm:"string;size:255" json:"-"`
	ApiToken  string `gogm:"string;size:255"  json:"-"`
	OpenId    string `gorm:"string;size:255"`
	GroupId   int64  `gorm:"group_id"`
}

func (user *User) AccessTo(admin contract.User) bool {
	return user.GetGroupId() == admin.GetGroupId()
}

func (user *User) GetGroupId() int64 {
	return user.GroupId
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

func (user *User) Login() (token string) {
	token = random.RandString(32)
	databases.Db.Model(user).Update("api_token", token)
	return
}
func (user *User) FindByName(username string) {
	databases.Db.Where("username= ?", username).Limit(1).First(user)
}
