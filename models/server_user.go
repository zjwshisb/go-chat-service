package models

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
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
	ApiToken string 	`gogm:"string;size:255"  json:"-"`
}

type ChatUser struct {
	ID        int64      `json:"id"`
	Username  string     `json:"username"`
	LastChatTime int64  `json:"last_chat_time"`
	Disabled bool `json:"disabled"`
	Online bool `json:"online"`
	Messages []interface{} `json:"messages"`
	Unread int `json:"unread"`
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
func (user *ServerUser) chatUsersKey() string {
	return fmt.Sprintf(serverChatUserKey, user.ID)
}
// 获取聊天过的用户
func (user *ServerUser) GetChatUsers() (users []*ChatUser) {
	ctx := context.Background()
	cmd := db.Redis.ZRevRangeWithScores(ctx, user.chatUsersKey(), 0, -1)
	if cmd.Err() != nil {
		return
	}
	uids := make([]int64, 0)
	for _, v := range cmd.Val() {
		member := v.Member.(string)
		id, err := strconv.ParseInt(member, 10, 64)
		if err == nil {
			uids = append(uids, id)
		}
	}
	if len(uids) == 0 {
		return
	}
	usesModel := make([]User, 0)
	db.Db.Find(&usesModel, uids)
	for _, v := range cmd.Val() {
		member := v.Member.(string)
		id, err := strconv.ParseInt(member, 10, 64)
		if err == nil {
			for _, u := range usesModel {
				if u.ID == id {
					item := &ChatUser{
						ID: id,
						Username: u.Username,
						LastChatTime: int64(v.Score),
						Messages: make([]interface{}, 0),
					}
					users = append(users, item)
					break
				}
			}
		}
	}
	return
}
// 移除聊天用户
func (user *ServerUser) RemoveChatUser(uid int64) error {
	ctx := context.Background()
	cmd := db.Redis.ZRem(ctx,  user.chatUsersKey(), uid)
	return cmd.Err()
}
// 更新聊天用户
func (user *ServerUser) UpdateChatUser(uid int64) error {
	ctx := context.Background()
	m := &redis.Z{Member: uid, Score: float64(time.Now().Unix())}
	cmd := db.Redis.ZAdd(ctx,  user.chatUsersKey(),  m)
	return cmd.Err()
}
