package models

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"time"
	"ws/internal/chat"
	"ws/internal/databases"
	"ws/internal/file"
	"ws/internal/util"
)

type BackendUser struct {
	ID        int64      `json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Username  string     `gorm:"string;size:255" json:"username"`
	Password  string     `gorm:"string;size:255" json:"-"`
	ApiToken string 	`gorm:"string;size:255"  json:"-"`
	Avatar string 		`gorm:"string;size:512" json:"-"`
	ShortcutReplies []*ShortcutReply `gorm:"foreignKey:UserId"`
}

func (user *BackendUser) GetPrimaryKey() int64 {
	return user.ID
}
func (user *BackendUser) GetAvatarUrl() string {
	if user.Avatar != "" {
		return file.Disk("local").Url(user.Avatar)
	}
	return ""
}
func (user *BackendUser) Login() (token string) {
	token = util.RandomStr(32)
	databases.Db.Model(user).Update("api_token", token)
	return
}
func (user *BackendUser) Logout()  {
	databases.Db.Model(user).Update("api_token", "")
}

func (user *BackendUser) Auth(c *gin.Context) bool {
	query := databases.Db.Where("api_token= ?", util.GetToken(c)).First(user)
	if query.Error == gorm.ErrRecordNotFound {
		return false
	}
	return true
}
func (user *BackendUser) FindByName(username string) bool {
	databases.Db.Where("username= ?", username).First(user)
	return user.ID > 0
}

func (user *BackendUser) GetTodayAcceptCount() (count int64) {
	ctx := context.Background()
	cmd := databases.Redis.ZRangeWithScores(ctx, chat.GetBackUserKey(user.GetPrimaryKey()), 0, -1)
	if cmd.Err() == redis.Nil {
		return
	}
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	timeNumber := t.Unix()
	for _, z := range cmd.Val() {
		score := int64(z.Score)
		if score > timeNumber {
			count ++
		}
	}
	return count
}