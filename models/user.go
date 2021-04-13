package models

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"strconv"
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

// 未发送给客服的消息
func (user *User) GetUnSendMsg() (messages []Message) {
	db.Db.Where("user_id = ?" , user.ID).Where("service_id", 0).Find(&messages)
	return
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
// 设置客服对象id
func (user *User) SetServerId(sid int64) error {
	if user.ID == 0 {
		return errors.New("user not exist")
	}
	ctx := context.Background()
	db.Redis.HSet(ctx, User2ServerHashKey, "test", sid)
	cmd := db.Redis.HSet(ctx, User2ServerHashKey, user.ID, sid)
	return cmd.Err()
}
// 获取最后一个客服id
func (user *User) GetLastServerId() int64 {
	if user.ID == 0 {
		return 0
	}
	ctx := context.Background()
	key := strconv.FormatInt(user.ID, 10)
	cmd := db.Redis.HGet(ctx, User2ServerHashKey, key)
	if sid, err := cmd.Int64(); err == nil {
		// 判断是否超时||已被客服端移除
		cmd := db.Redis.ZScore(ctx, fmt.Sprintf(serverChatUserKey, sid), key)
		if cmd.Err() == redis.Nil {
			return 0
		}
		t := int64(cmd.Val())
		if t <= (time.Now().Unix() - 3600 * 24 * 2) {
			return 0
		}
		return sid
	}
	return 0
}