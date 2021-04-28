package models

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
	"ws/config"
	"ws/core/image"
	"ws/db"
	"ws/util"
)

const (
	serverChatUserKey = "server-user:%d:chat-user"
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
	ApiToken string 	`gorm:"string;size:255"  json:"-"`
	Avatar string 		`gorm:"string;size:512" json:"-"`
}



func (user *ServerUser) GetAvatarUrl() string {
	if user.Avatar != "" {
		return image.Url(user.Avatar)
	}
	return ""
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
	db.Db.Where("api_token= ?", util.GetToken(c)).First(user)
}

func (user *ServerUser) FindByName(username string) () {
	db.Db.Where("username= ?", username).First(user)
}
func (user *ServerUser) ChatUsersKey() string {
	return fmt.Sprintf(serverChatUserKey, user.ID)
}
// 检查聊天对象是否过期
func (user *ServerUser) CheckChatUserLegal(uid int64) bool {
	ctx := context.Background()
	cmd := db.Redis.ZScore(ctx, user.ChatUsersKey(), strconv.FormatInt(uid , 10))
	if cmd.Err() == redis.Nil {
		return false
	}
	score := cmd.Val()
	t := int64(score)
	if (time.Now().Unix() - t) <= config.App.ChatSessionDuration * 24 * 60 * 60 {
		return true
	}
	return false
}
// 获取聊天过的用户
func (user *ServerUser) GetChatUsers() (users []*User) {
	users = make([]*User, 0)
	ctx := context.Background()
	cmd := db.Redis.ZRange(ctx, user.ChatUsersKey(), 0, -1)
	if cmd.Err() == redis.Nil {
		return
	}
	uids := make([]int64, 0)
	for _, idStr := range cmd.Val() {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err == nil {
			uids = append(uids, id)
		}
	}
	if len(uids) == 0 {
		return
	}
	db.Db.Find(&users, uids)
	fmt.Println(users)
	return
}
// 移除聊天用户
func (user *ServerUser) RemoveChatUser(uid int64) error {
	ctx := context.Background()
	cmd := db.Redis.ZRem(ctx,  user.ChatUsersKey(), uid)
	return cmd.Err()
}
// 更新聊天用户
func (user *ServerUser) UpdateChatUser(uid int64) error {
	ctx := context.Background()
	m := &redis.Z{Member: uid, Score: float64(time.Now().Unix())}
	cmd := db.Redis.ZAdd(ctx,  user.ChatUsersKey(),  m)
	return cmd.Err()
}
