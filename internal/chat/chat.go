package chat

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
	"time"
	"ws/configs"
	"ws/internal/databases"
	"ws/internal/util"
)


const (
	// 用户 => 客服 hashes
	user2ServerHashKey = "user-to-server"
	// 客服 => {value: userId, source: limitTime}[] sorted sets
	serverChatUserKey = "server-user:%d:chat-user"
	// 客服 => {uid: lastTime} hashes
	serverUserLastChatKey = "server-user:%d:chat-user:last-time"
	// 待人工接入的用户 sets
	manualUserKey = "user:manual"
)
// 系统头像
func SystemAvatar() string  {
	return  util.PublicAsset("avatar.jpeg")
}
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

// 转接人工客服的用户
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
// 获取聊天过的用户ids以及对应的最后聊天时间
func GetChatUserIds(sid int64)  ([]int64, []int64) {
	ctx := context.Background()
	cmd := databases.Redis.ZRangeWithScores(ctx, GetBackUserKey(sid), 0, -1)
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
// 设置客服用户最后聊天时间
func SetServerUserLastChatTime(uid int64,sid int64) error {
	ctx := context.Background()
	cmd := databases.Redis.HSet(ctx, fmt.Sprintf(serverUserLastChatKey, sid), uid, time.Now().Unix())
	return cmd.Err()
}
// 获取客服用户最后聊天时间
func GetServerUserLastChatTime(uid int64, sid int64)  int64 {
	ctx := context.Background()
	cmd := databases.Redis.HGet(ctx, fmt.Sprintf(serverUserLastChatKey, sid), strconv.FormatInt(uid, 10))
	t, _ := strconv.ParseInt(cmd.Val(), 10, 64)
	return t
}
// 设置用户客服对象id
func SetUserServerId(uid int64,sid int64, duration int64) error {
	ctx := context.Background()
	cmd := databases.Redis.HSet(ctx, user2ServerHashKey,uid, sid)
	_ = UpdateUserServerId(uid, sid, duration)
	_ = RemoveManual(uid)
	return cmd.Err()
}
// 更新会话时间
func UpdateUserServerId(uid int64, sid int64, duration int64) error {
	ctx := context.Background()
	m := &redis.Z{Member: uid, Score: float64(time.Now().Unix() + duration)}
	_ = SetServerUserLastChatTime(uid, sid)
	cmd := databases.Redis.ZAdd(ctx, GetBackUserKey(sid),  m)
	return cmd.Err()
}
// 清除用户客服对象id
func RemoveUserServerId(uid int64, sid int64) error {
	ctx := context.Background()
	cmd := databases.Redis.HDel(ctx, user2ServerHashKey, strconv.FormatInt(uid, 10))
	cmd = databases.Redis.HDel(ctx, fmt.Sprintf(serverUserLastChatKey, sid), strconv.FormatInt(uid, 10))
	cmd = databases.Redis.ZRem(ctx, GetBackUserKey(sid), uid)
	return cmd.Err()
}

// 获取用户最后一个会话客服id
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
		if t <=  time.Now().Unix() {
			return 0
		}
		return sid
	}
	return 0
}
// 客服给用户发消息的会话有效期, 既用户在这时间内可以回复客服
func GetUserSessionSecond() int64 {
	setting := Settings[UserSessionDuration]
	dayFloat, err := strconv.ParseFloat(setting.GetValue(), 64)
	if err != nil {
		log.Fatal(err)
	}
	second := int64(dayFloat* 24 * 60 * 60)
	return second
}
// 用户给客服发消息的会话有效期, 既客服在这时间内可以回复用户
func GetServiceSessionSecond() int64 {
	setting := Settings[ServiceSessionDuration]
	dayFloat, err := strconv.ParseFloat(setting.GetValue(), 64)
	if err != nil {
		log.Fatal(err)
	}
	second := int64(dayFloat * 24 * 60 * 60)
	return second
}
// 客服的聊天用户SortedSet 的key
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
	limitTime := int64(score)
	return limitTime > time.Now().Unix()
}
// 标记 用户微信订阅消息 已订阅
func SetSubscribe(uid int64) error {
	ctx := context.Background()
	templateId := configs.Wechat.SubscribeTemplateIdOne
	key := fmt.Sprintf("user:%d:subscribe:%s", uid, templateId)
	cmd := databases.Redis.Set(ctx, key, 1, 0)
	return cmd.Err()
}
// 查询 用户微信订阅消息
func IsSubScribe(uid int64) bool  {
	ctx := context.Background()
	templateId := configs.Wechat.SubscribeTemplateIdOne
	key := fmt.Sprintf("user:%d:subscribe:%s", uid, templateId)
	cmd := databases.Redis.Get(ctx, key)
	if cmd.Err() == redis.Nil {
		return false
	}
	return true
}
// 删除 用户微信订阅消息 标记
func DelSubScribe(uid int64) bool {
	ctx := context.Background()
	templateId := configs.Wechat.SubscribeTemplateIdOne
	key := fmt.Sprintf("user:%d:subscribe:%s", uid, templateId)
	databases.Redis.Del(ctx, key)
	return true
}