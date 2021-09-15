package chat

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
	"time"
	"ws/app/databases"
)

const (
	// 用户 => 客服 hashes
	user2AdminHashKey = "user-to-admin"
	// 客服 => {value: userId, source: limitTime}[] sorted sets
	adminChatUserKey = "admin:%d:chat-user"
	// 客服 => {uid: lastTime} hashes
	adminUserLastChatKey = "admin:%d:chat-user:last-time"
	// 待人工接入的用户 sets
	manualUserKey = "user:manual"
)

// 添加用户到人工客服列表
func AddToManual(uid int64) error  {
	ctx := context.Background()
	cmd := databases.Redis.SAdd(ctx, manualUserKey, uid)
	return cmd.Err()
}
// 判断用户是否在人工客服等待区
func IsInManual(uid int64) bool  {
	ctx := context.Background()
	cmd := databases.Redis.SIsMember(ctx, manualUserKey, uid)
	return cmd.Val()
}
// 从人工客服列表移除用户
func RemoveManual(uid int64) error {
	ctx := context.Background()
	cmd := databases.Redis.SRem(ctx, manualUserKey, uid)
	return cmd.Err()
}
// 待接入的用户ids(待人工接待列表)
func GetManualUserIds() []int64 {
	ctx := context.Background()
	cmd := databases.Redis.SMembers(ctx, manualUserKey)
	uid := make([]int64, 0, len(cmd.Val()))
	for _, uidStr := range cmd.Val() {
		id , err := strconv.ParseInt(uidStr, 10, 64)
		if err == nil {
			uid = append(uid, id)
		}
	}
	return uid
}
// 获取聊天过的用户ids以及对应的有效时间
func GetAdminUserIds(adminId int64)  ([]int64, []int64) {
	ctx := context.Background()
	cmd := databases.Redis.ZRevRangeWithScores(ctx, GetAdminUserKey(adminId), 0, -1)
	uids := make([]int64, 0, len(cmd.Val()))
	times :=  make([]int64, 0, len(cmd.Val()))
	for _, item := range cmd.Val() {
		id, err := strconv.ParseInt(item.Member.(string), 10, 64)
		if err == nil {
			uids = append(uids, id)
		}
		score := int64(item.Score)
		times = append(times, score)
	}
	return uids, times
}
func GetAdminUserActiveCount(adminId int64) int {
	ctx := context.Background()
	cmd := databases.Redis.ZRangeByScore(ctx, GetAdminUserKey(adminId), &redis.ZRangeBy{
		Min:    strconv.FormatInt(time.Now().Unix(), 10),
		Max:    "+inf",
	})
	return len(cmd.Val())
}
// 设置客服用户最后聊天时间
func SetAdminUserLastChatTime(uid int64,adminId int64) error {
	ctx := context.Background()
	cmd := databases.Redis.HSet(ctx, fmt.Sprintf(adminUserLastChatKey, adminId), uid, time.Now().Unix())
	return cmd.Err()
}
// 获取客服用户最后聊天时间
func GetAdminUserLastChatTime(uid int64, adminId int64)  int64 {
	ctx := context.Background()
	cmd := databases.Redis.HGet(ctx, fmt.Sprintf(adminUserLastChatKey, adminId), strconv.FormatInt(uid, 10))
	t, _ := strconv.ParseInt(cmd.Val(), 10, 64)
	return t
}
// 设置用户客服对象id
func SetUserAdminId(uid int64,adminId int64, duration int64) error {
	ctx := context.Background()
	cmd := databases.Redis.HSet(ctx, user2AdminHashKey, uid, adminId)
	_ = UpdateUserAdminId(uid, adminId, duration)
	_ = RemoveManual(uid)
	return cmd.Err()
}
// 更新会话时间
func UpdateUserAdminId(uid int64, adminId int64, duration int64) error {
	ctx := context.Background()
	m := &redis.Z{Member: uid, Score: float64(time.Now().Unix() + duration)}
	_ = SetAdminUserLastChatTime(uid, adminId)
	cmd := databases.Redis.ZAdd(ctx, GetAdminUserKey(adminId),  m)
	return cmd.Err()
}
// 清除用户客服对象id
func RemoveUserAdminId(uid int64, adminId int64) error {
	ctx := context.Background()
	cmd := databases.Redis.HDel(ctx, user2AdminHashKey, strconv.FormatInt(uid, 10))
	cmd = databases.Redis.HDel(ctx, fmt.Sprintf(adminUserLastChatKey, adminId), strconv.FormatInt(uid, 10))
	cmd = databases.Redis.ZRem(ctx, GetAdminUserKey(adminId), uid)
	return cmd.Err()
}

// 获取用户最后一个会话客服id
func GetUserLastAdminId(uid int64) int64 {
	ctx := context.Background()
	key := strconv.FormatInt(uid, 10)
	cmd := databases.Redis.HGet(ctx, user2AdminHashKey, key)
	if sid, err := cmd.Int64(); err == nil {
		// 判断是否超时|已被客服移除
		cmd := databases.Redis.ZScore(ctx, GetAdminUserKey(sid), key)
		if cmd.Err() == redis.Nil {
			return 0
		}
		t := int64(cmd.Val())
		if t <=  time.Now().Unix() {
			return 0
		}
		return sid
	}
	return 0
}

// 客服的聊天用户SortedSet 的key
func GetAdminUserKey(adminId int64) string {
	return fmt.Sprintf(adminChatUserKey, adminId)
}

// 检查用户对于客服是否合法
func CheckUserIdLegal(uid int64, adminId int64) bool {
	ctx := context.Background()
	cmd := databases.Redis.ZScore(ctx, GetAdminUserKey(adminId), strconv.FormatInt(uid , 10))
	if cmd.Err() == redis.Nil {
		return false
	}
	score := cmd.Val()
	limitTime := int64(score)
	return limitTime > time.Now().Unix()
}
// 离线时超过多久就自动断开会话
func GetOfflineDuration() int64 {
	setting := Settings[MinuteToBreak]
	minuteStr := setting.GetValue()
	minute, err := strconv.ParseInt(minuteStr, 10,64)
	if err != nil {
		log.Fatal(err)
	}
	return minute * 60
}
