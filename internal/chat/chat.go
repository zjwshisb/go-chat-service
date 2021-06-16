package chat

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
	"ws/configs"
	"ws/internal/databases"
)

const (
	user2ServerHashKey = "user-to-server"
	serverChatUserKey = "server-user:%d:chat-user"
)

// 设置用户客服对象id
func SetUserServerId(uid int64,sid int64) error {
	ctx := context.Background()
	cmd := databases.Redis.HSet(ctx, user2ServerHashKey,uid, sid)
	_ = UpdateUserServerId(uid, sid)
	return cmd.Err()
}
// 更新客服的用户的最后聊天时间
func UpdateUserServerId(uid int64, sid int64) error {
	ctx := context.Background()
	m := &redis.Z{Member: uid, Score: float64(time.Now().Unix())}
	cmd := databases.Redis.ZAdd(ctx, GetBackUserKey(sid),  m)
	return cmd.Err()
}
// 清除用户客服对象id
func RemoveUserServerId(uid int64, sid int64) error {
	ctx := context.Background()
	cmd := databases.Redis.HDel(ctx, user2ServerHashKey, strconv.FormatInt(uid, 10))
	if cmd.Err() != nil {
		return cmd.Err()
	}
	cmd = databases.Redis.ZRem(ctx, GetBackUserKey(sid), uid)
	return cmd.Err()
}
// 获取用户最后一个客服id
func GetUserLastServerId(uid int64) int64 {
	ctx := context.Background()
	key := strconv.FormatInt(uid, 10)
	cmd := databases.Redis.HGet(ctx, user2ServerHashKey, key)
	if sid, err := cmd.Int64(); err == nil {
		// 判断是否超时|已被客服移除
		cmd := databases.Redis.ZScore(ctx, GetBackUserKey(sid), key)
		if cmd.Err() == redis.Nil {
			return 0
		}
		t := int64(cmd.Val())
		if t <= (time.Now().Unix() - configs.App.ChatSessionDuration * 24 * 60 * 60) {
			return 0
		}
		return sid
	}
	return 0
}
// 获取redis 客服的聊天用户SortedSet 的key
func GetBackUserKey(sid int64) string {
	return fmt.Sprintf(serverChatUserKey, sid)
}
// 检查用户对于客服是否合法
func CheckUserIdLegal(uid int64, sid int64) bool {
	ctx := context.Background()
	cmd := databases.Redis.ZScore(ctx, GetBackUserKey(sid), strconv.FormatInt(uid , 10))
	if cmd.Err() == redis.Nil {
		return false
	}
	score := cmd.Val()
	t := int64(score)
	if (time.Now().Unix() - t) <= configs.App.ChatSessionDuration * 24 * 60 * 60 {
		return true
	}
	return false
}
