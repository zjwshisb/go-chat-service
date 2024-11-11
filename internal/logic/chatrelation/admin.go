package chat

import (
	"context"
	"errors"
	"fmt"
	"gf-chat/internal/service"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gogf/gf/v2/frame/g"
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

type ctx = context.Context

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
func (s *sChatRelation) AddUser(ctx ctx, adminId uint, uid uint) error {
	_ = s.SetUserAdmin(ctx, uid, adminId)
	g.Redis().Do(ctx, "ZAdd", s.getUserCacheKey(adminId), time.Now().Unix()+DefaultSessionTime, uid)
	err := s.UpdateUser(ctx, adminId, uid)
	return err
}

// UpdateUser 更新user
// 更新limit time
// 更新最后聊天时间
func (s *sChatRelation) UpdateUser(ctx ctx, adminId uint, uid uint) error {
	err := s.UpdateLimitTime(ctx, adminId, uid, DefaultSessionTime)
	if err != nil {
		return err
	}
	err = s.UpdateLastChatTime(ctx, adminId, uid)
	return err
}

// RemoveUser 移除user
func (s *sChatRelation) RemoveUser(ctx ctx, adminId uint, uid uint) error {
	_ = s.RemoveUserAdmin(ctx, uid)
	_ = s.RemoveLastChatTime(ctx, adminId, uid)
	_, err := g.Redis().Do(ctx, "Zrem", s.getUserCacheKey(adminId), uid)
	return err
}

// IsUserValid 检查用户对于客服是否合法
func (s *sChatRelation) IsUserValid(ctx ctx, adminId uint, uid uint) bool {
	b := s.GetLimitTime(ctx, adminId, uid) > time.Now().Unix()
	return b
}

// IsUserExist user是否存在
func (s *sChatRelation) IsUserExist(ctx ctx, adminId uint, uid uint) bool {
	_, err := g.Redis().Do(ctx, "ZScore", s.getUserCacheKey(adminId), uid)
	if err == redis.Nil {
		return false
	}
	return true
}

// GetLastChatTime 获取最后聊天时间
func (s *sChatRelation) GetLastChatTime(ctx ctx, adminId uint, uid uint) uint {
	val, _ := g.Redis().Do(ctx, "HGet", fmt.Sprintf(adminUserLastChatKey, adminId), uid)
	return val.Uint()
}

// RemoveLastChatTime 移除最后聊天时间
func (s *sChatRelation) RemoveLastChatTime(ctx ctx, adminId uint, uid uint) error {
	_, err := g.Redis().Do(ctx, "HDel", fmt.Sprintf(adminUserLastChatKey, adminId), uid)
	return err
}

// UpdateLastChatTime 更新最后聊天时间
func (s *sChatRelation) UpdateLastChatTime(ctx ctx, adminId uint, uid uint) error {
	_, err := g.Redis().Do(ctx, "HSet", fmt.Sprintf(adminUserLastChatKey, adminId), uid, time.Now().Unix())
	return err
}

// GetActiveCount 获取有效的用户数量
func (s *sChatRelation) GetActiveCount(ctx ctx, adminId uint) uint {
	val, _ := g.Redis().Do(ctx, "ZRangeByScore",
		s.getUserCacheKey(adminId), time.Now().Unix(), "+inf")
	return uint(len(val.Int64s()))
}

// UpdateLimitTime 更新有效期
func (s *sChatRelation) UpdateLimitTime(ctx ctx, adminId uint, uid uint, duration int64) error {
	if !s.IsUserExist(ctx, adminId, uid) {
		return errors.New("user not valid")
	}
	_, err := g.Redis().Do(ctx, "ZAdd", s.getUserCacheKey(adminId), time.Now().Unix()+duration, uid)
	return err
}

// GetLimitTime 获取有效期
func (s *sChatRelation) GetLimitTime(ctx ctx, adminId uint, uid uint) int64 {
	val, _ := g.Redis().Do(ctx, "ZScore", s.getUserCacheKey(adminId), uid)
	return val.Int64()
}

func (s *sChatRelation) GetInvalidUsers(ctx ctx, adminId uint) []uint {
	val, _ := g.Redis().Do(ctx, "zrangebyscore",
		s.getUserCacheKey(adminId), "-inf", time.Now().Unix())
	return val.Uints()
}

// GetUsersWithLimitTime 获取所有user以及对应的有效期
func (s *sChatRelation) GetUsersWithLimitTime(ctx ctx, adminId uint) ([]uint, []int64) {
	val, _ := g.Redis().Do(ctx, "ZREVRANGE", s.getUserCacheKey(adminId), 0, -1, "WITHSCORES")
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
func (s *sChatRelation) SetUserAdmin(ctx ctx, uid uint, adminId uint) error {
	_, err := g.Redis().Do(ctx, "hset", user2AdminHashKey, uid, adminId)
	return err
}

// RemoveUserAdmin RemoveAdmin 移除用户客服
func (s *sChatRelation) RemoveUserAdmin(ctx ctx, uid uint) error {
	_, err := g.Redis().Do(ctx, "hdel", user2AdminHashKey, uid)
	return err
}

// GetUserValidAdmin GetValidAdmin 获取用户客服
func (s *sChatRelation) GetUserValidAdmin(ctx ctx, uid uint) uint {
	val, err := g.Redis().Do(ctx, "HGet", user2AdminHashKey, uid)
	if err == nil {
		adminId := val.Uint()
		limitTime := s.GetLimitTime(ctx, adminId, uid)
		if limitTime > time.Now().Unix() {
			return val.Uint()
		}
		// 无效了直接清除掉
		_ = s.RemoveUserAdmin(ctx, uid)
	}
	return 0
}
