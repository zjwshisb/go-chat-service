package chat

import (
	"context"
	"errors"
	"fmt"
	"gf-chat/internal/service"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
)

const (
	// 客服 => {value: userId, source: limitTime}[] sorted sets
	adminChatUserKey = "admin:%d:chat-user"
	// 客服 => {uid: lastTime} hashes
	adminUserLastChatKey = "admin:%d:chat-user:last-time"

	// 用户 => 客服 hashes
	user2AdminHashKey = "user-to-admin"
)

var (
	DefaultSessionTime int64 = 24 * 60 * 60
)

func init() {
	service.RegisterChatRelation(&sChatRelation{})
}

type sChatRelation struct {
}

func (s *sChatRelation) getUserCacheKey(adminId uint) string {
	return fmt.Sprintf(adminChatUserKey, adminId)
}

// AddUser 接入user
func (s *sChatRelation) AddUser(adminId, uid uint) error {
	ctx := context.Background()
	_ = s.SetUserAdmin(uid, adminId)
	g.Redis().Do(ctx, "ZAdd", s.getUserCacheKey(adminId), time.Now().Unix()+DefaultSessionTime, uid)
	err := s.UpdateUser(adminId, uid)
	return err
}

// UpdateUser 更新user
// 更新limit time
// 更新最后聊天时间
func (s *sChatRelation) UpdateUser(adminId uint, uid uint) error {
	err := s.UpdateLimitTime(adminId, uid, DefaultSessionTime)
	if err != nil {
		return err
	}
	err = s.UpdateLastChatTime(adminId, uid)
	return err
}

// RemoveUser 移除user
func (s *sChatRelation) RemoveUser(adminId uint, uid uint) error {
	ctx := context.Background()
	_ = s.RemoveUserAdmin(uid)
	_ = s.RemoveLastChatTime(adminId, uid)
	_, err := g.Redis().Do(ctx, "Zrem", s.getUserCacheKey(adminId), uid)
	return err
}

// IsUserValid 检查用户对于客服是否合法
func (s *sChatRelation) IsUserValid(adminId uint, uid uint) bool {
	b := s.GetLimitTime(adminId, uid) > time.Now().Unix()
	return b
}

// IsUserExist user是否存在
func (s *sChatRelation) IsUserExist(adminId uint, uid uint) bool {
	ctx := context.Background()
	_, err := g.Redis().Do(ctx, "ZScore", s.getUserCacheKey(adminId), uid)
	if err == redis.Nil {
		return false
	}
	return true
}

// GetLastChatTime 获取最后聊天时间
func (s *sChatRelation) GetLastChatTime(adminId uint, uid uint) int64 {
	ctx := context.Background()
	val, _ := g.Redis().Do(ctx, "HGet", fmt.Sprintf(adminUserLastChatKey, adminId), uid)
	return val.Int64()
}

// RemoveLastChatTime 移除最后聊天时间
func (s *sChatRelation) RemoveLastChatTime(adminId uint, uid uint) error {
	ctx := context.Background()
	_, err := g.Redis().Do(ctx, "HDel", fmt.Sprintf(adminUserLastChatKey, adminId), uid)
	return err
}

// UpdateLastChatTime 更新最后聊天时间
func (s *sChatRelation) UpdateLastChatTime(adminId uint, uid uint) error {
	_, err := g.Redis().Do(gctx.New(), "HSet", fmt.Sprintf(adminUserLastChatKey, adminId), uid, time.Now().Unix())
	return err
}

// GetActiveCount 获取有效的用户数量
func (s *sChatRelation) GetActiveCount(adminId uint) uint {
	val, _ := g.Redis().Do(gctx.New(), "ZRangeByScore",
		s.getUserCacheKey(adminId), time.Now().Unix(), "+inf")
	return uint(len(val.Int64s()))
}

// UpdateLimitTime 更新有效期
func (s *sChatRelation) UpdateLimitTime(adminId uint, uid uint, duration int64) error {
	if !s.IsUserExist(adminId, uid) {
		return errors.New("user not valid")
	}
	_, err := g.Redis().Do(gctx.New(), "ZAdd", s.getUserCacheKey(adminId), time.Now().Unix()+duration, uid)
	return err
}

// GetLimitTime 获取有效期
func (s *sChatRelation) GetLimitTime(adminId uint, uid uint) int64 {
	val, _ := g.Redis().Do(gctx.New(), "ZScore", s.getUserCacheKey(adminId), uid)
	return val.Int64()
}

func (s *sChatRelation) GetInvalidUsers(adminId uint) []uint {
	val, _ := g.Redis().Do(gctx.New(), "zrangebyscore",
		s.getUserCacheKey(adminId), "-inf", time.Now().Unix())
	return val.Uints()
}

// GetUsersWithLimitTime 获取所有user以及对应的有效期
func (s *sChatRelation) GetUsersWithLimitTime(adminId uint) ([]uint, []int64) {
	val, _ := g.Redis().Do(gctx.New(), "ZREVRANGE", s.getUserCacheKey(adminId), 0, -1, "WITHSCORES")
	uids := make([]uint, 0, len(val.Slice()))
	times := make([]int64, 0, len(val.Slice()))
	for index, item := range val.Slice() {
		types := index % 2
		switch types {
		case 0:
			uids = append(uids, gconv.Uint(item))
		case 1:
			times = append(times, gconv.Int64(item))
		}
	}
	return uids, times
}

// SetUserAdmin SetAdmin 设置用户客服
func (s *sChatRelation) SetUserAdmin(uid uint, adminId uint) error {
	_, err := g.Redis().Do(gctx.New(), "hset", user2AdminHashKey, uid, adminId)
	return err
}

// RemoveUserAdmin RemoveAdmin 移除用户客服
func (s *sChatRelation) RemoveUserAdmin(uid uint) error {
	_, err := g.Redis().Do(gctx.New(), "hdel", user2AdminHashKey, uid)
	return err
}

// GetUserValidAdmin GetValidAdmin 获取用户客服
func (s *sChatRelation) GetUserValidAdmin(uid uint) uint {
	val, err := g.Redis().Do(gctx.New(), "HGet", user2AdminHashKey, uid)
	if err == nil {
		adminId := val.Uint()
		limitTime := s.GetLimitTime(adminId, uid)
		if limitTime > time.Now().Unix() {
			return val.Uint()
		}
		// 无效了直接清除掉
		_ = s.RemoveUserAdmin(uid)
	}
	return 0
}
