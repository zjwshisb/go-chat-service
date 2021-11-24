package chat

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
	"ws/app/databases"
)
const (
	// 客服 => {value: userId, source: limitTime}[] sorted sets
	adminChatUserKey = "admin:%d:chat-user"
	// 客服 => {uid: lastTime} hashes
	adminUserLastChatKey = "admin:%d:chat-user:last-time"
)
var (
	AdminService = &adminService{}
)

type adminService struct {
}
// cache key
func (adminService *adminService) getUserCacheKey(adminId int64) string  {
	return fmt.Sprintf(adminChatUserKey, adminId)
}
// 接入user
func (adminService *adminService) AddUser(adminId int64, uid int64, duration int64) error  {
	ctx := context.Background()
	_ = UserService.SetAdmin(uid, adminId)
	m := &redis.Z{Member: uid, Score: float64(time.Now().Unix() + duration)}
	_ = databases.Redis.ZAdd(ctx, AdminService.getUserCacheKey(adminId),  m)
	err := ManualService.Remove(uid)
	return err
}
// 更新user
// 更新limit time
// 更新最后聊天时间
func (adminService *adminService) UpdateUser(adminId int64, uid int64, duration int64) error {
	err := adminService.UpdateLimitTime(adminId, uid, duration)
	if err != nil {
		return err
	}
	err = adminService.UpdateLastChatTime(adminId, uid)
	return err
}
// 移除user
func (adminService *adminService) RemoveUser(adminId int64, uid int64) error  {
	ctx := context.Background()
	_ = UserService.RemoveAdmin(uid)
	_ = adminService.RemoveLastChatTime(adminId, uid)
	cmd := databases.Redis.ZRem(ctx, AdminService.getUserCacheKey(adminId), uid)
	return cmd.Err()
}

// 检查用户对于客服是否合法
func (adminService *adminService) IsUserValid(adminId int64, uid int64) bool {
	b := adminService.GetLimitTime(adminId, uid) > time.Now().Unix()
	return b
}

// user是否存在
func (adminService *adminService) IsUserExist(adminId int64, uid int64) bool  {
	ctx := context.Background()
	exist := databases.Redis.ZScore(ctx, adminService.getUserCacheKey(adminId), strconv.FormatInt(uid , 10))
	if exist.Err() == redis.Nil  {
		return false
	}
	return true
}
// 获取最后聊天时间
func (adminService *adminService) GetLastChatTime(adminId int64, uid int64) int64 {
	ctx := context.Background()
	cmd := databases.Redis.HGet(ctx, fmt.Sprintf(adminUserLastChatKey, adminId), strconv.FormatInt(uid, 10))
	t, _ := strconv.ParseInt(cmd.Val(), 10, 64)
	return t
}
// 移除最后聊天时间
func (adminService *adminService) RemoveLastChatTime(adminId int64, uid int64) error {
	ctx := context.Background()
	cmd := databases.Redis.HDel(ctx, fmt.Sprintf(adminUserLastChatKey, adminId), strconv.FormatInt(uid, 10))
	return cmd.Err()
}
// 更新最后聊天时间
func (adminService *adminService) UpdateLastChatTime(adminId int64, uid int64) error  {
	ctx := context.Background()
	cmd := databases.Redis.HSet(ctx, fmt.Sprintf(adminUserLastChatKey, adminId), uid, time.Now().Unix())
	return cmd.Err()
}
// 获取有效的用户数量
func (adminService *adminService) GetActiveCount(adminId int64) int  {
	ctx := context.Background()
	cmd := databases.Redis.ZRangeByScore(ctx, adminService.getUserCacheKey(adminId), &redis.ZRangeBy{
		Min:    strconv.FormatInt(time.Now().Unix(), 10),
		Max:    "+inf",
	})
	return len(cmd.Val())
}
// 更新有效期
func (adminService *adminService) UpdateLimitTime(adminId int64, uid int64, duration int64) error {
	if !adminService.IsUserExist(adminId, uid)  {
		return errors.New("user not valid")
	}
	ctx := context.Background()
	m := &redis.Z{Member: uid, Score: float64(time.Now().Unix() + duration)}
	cmd1 := databases.Redis.ZAdd(ctx, AdminService.getUserCacheKey(adminId),  m)
	return cmd1.Err()
}
// 获取有效期
func (adminService *adminService) GetLimitTime(adminId int64, uid int64) int64 {
	ctx := context.Background()
	cmd := databases.Redis.ZScore(ctx, adminService.getUserCacheKey(adminId), strconv.FormatInt(uid , 10))
	if cmd.Err() == redis.Nil {
		return 0
	}
	score := cmd.Val()
	return int64(score)
}
// 获取所有user
func (adminService *adminService) GetUsersWithLimitTime(adminId int64) ([]int64, []int64)  {
	ctx := context.Background()
	cmd := databases.Redis.ZRevRangeWithScores(ctx, adminService.getUserCacheKey(adminId), 0, -1)
	uids := make([]int64, 0, len(cmd.Val()))
	times :=  make([]int64, 0, len(cmd.Val()))
	for _, item := range cmd.Val() {
		id, err := strconv.ParseInt(item.Member.(string), 10, 64)
		if err == nil {
			uids = append(uids, id)
			score := int64(item.Score)
			times = append(times, score)
		}
	}
	return uids, times
}
