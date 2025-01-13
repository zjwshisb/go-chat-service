package chatrelation

import (
	"errors"
	"fmt"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
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
func (s *sChatRelation) AddUser(ctx gctx.Ctx, adminId uint, uid uint) (err error) {
	_ = s.SetUserAdmin(ctx, uid, adminId)
	_, err = g.Redis().Do(ctx, "ZAdd", s.getUserCacheKey(adminId), time.Now().Unix()+DefaultSessionTime, uid)
	if err != nil {
		return
	}
	err = s.UpdateUser(ctx, adminId, uid)
	return err
}

// UpdateUser 更新user
// 更新limit time
// 更新最后聊天时间
func (s *sChatRelation) UpdateUser(ctx gctx.Ctx, adminId uint, uid uint) (err error) {
	err = s.UpdateLimitTime(ctx, adminId, uid, DefaultSessionTime)
	if err != nil {
		return err
	}
	err = s.UpdateLastChatTime(ctx, adminId, uid)
	return err
}

// RemoveUser 移除user
func (s *sChatRelation) RemoveUser(ctx gctx.Ctx, adminId uint, uid uint) (err error) {
	err = s.RemoveUserAdmin(ctx, uid)
	if err != nil {
		return err
	}
	err = s.RemoveLastChatTime(ctx, adminId, uid)
	if err != nil {
		return err
	}
	_, err = g.Redis().Do(ctx, "Zrem", s.getUserCacheKey(adminId), uid)
	return err
}

// IsUserValid 检查用户对于客服是否合法
func (s *sChatRelation) IsUserValid(ctx gctx.Ctx, adminId uint, uid uint) bool {
	b := s.GetLimitTime(ctx, adminId, uid) > time.Now().Unix()
	return b
}

// IsUserExist user是否存在
func (s *sChatRelation) IsUserExist(ctx gctx.Ctx, adminId uint, uid uint) bool {
	_, err := g.Redis().Do(ctx, "ZScore", s.getUserCacheKey(adminId), uid)
	if errors.Is(err, redis.Nil) {
		return false
	}
	return true
}

// GetLastChatTime 获取最后聊天时间
func (s *sChatRelation) GetLastChatTime(ctx gctx.Ctx, adminId uint, uid uint) uint {
	val, _ := g.Redis().Do(ctx, "HGet", fmt.Sprintf(adminUserLastChatKey, adminId), uid)
	return val.Uint()
}

// RemoveLastChatTime 移除最后聊天时间
func (s *sChatRelation) RemoveLastChatTime(ctx gctx.Ctx, adminId uint, uid uint) error {
	_, err := g.Redis().Do(ctx, "HDel", fmt.Sprintf(adminUserLastChatKey, adminId), uid)
	return err
}

// UpdateLastChatTime 更新最后聊天时间
func (s *sChatRelation) UpdateLastChatTime(ctx gctx.Ctx, adminId uint, uid uint) error {
	_, err := g.Redis().Do(ctx, "HSet", fmt.Sprintf(adminUserLastChatKey, adminId), uid, time.Now().Unix())
	return err
}

// GetActiveCount 获取有效的用户数量
func (s *sChatRelation) GetActiveCount(ctx gctx.Ctx, adminId uint) uint {
	val, _ := g.Redis().Do(ctx, "ZRangeByScore",
		s.getUserCacheKey(adminId), time.Now().Unix(), "+inf")
	return uint(len(val.Int64s()))
}

// UpdateLimitTime 更新有效期
func (s *sChatRelation) UpdateLimitTime(ctx gctx.Ctx, adminId uint, uid uint, duration int64) error {
	if !s.IsUserExist(ctx, adminId, uid) {
		return gerror.New("user not valid")
	}
	_, err := g.Redis().Do(ctx, "ZAdd", s.getUserCacheKey(adminId), time.Now().Unix()+duration, uid)
	return err
}

// GetLimitTime 获取有效期
func (s *sChatRelation) GetLimitTime(ctx gctx.Ctx, adminId uint, uid uint) int64 {
	val, _ := g.Redis().Do(ctx, "ZScore", s.getUserCacheKey(adminId), uid)
	return val.Int64()
}

func (s *sChatRelation) GetInvalidUsers(ctx gctx.Ctx, adminId uint) []uint {
	val, _ := g.Redis().Do(ctx, "zrangebyscore",
		s.getUserCacheKey(adminId), "-inf", time.Now().Unix())
	return val.Uints()
}

// GetUsersWithLimitTime 获取所有user以及对应的有效期
func (s *sChatRelation) GetUsersWithLimitTime(ctx gctx.Ctx, adminId uint) (uids []uint, times []int64) {
	val, _ := g.Redis().Do(ctx, "ZREVRANGE", s.getUserCacheKey(adminId), 0, -1, "WITHSCORES")
	for _, item := range val.Vars() {
		s := item.Slice()
		uids = append(uids, gconv.Uint(s[0]))
		times = append(times, gconv.Int64(s[1]))
	}
	return
}

// SetUserAdmin SetAdmin 设置用户客服
func (s *sChatRelation) SetUserAdmin(ctx gctx.Ctx, uid uint, adminId uint) (err error) {
	_, err = g.Redis().Do(ctx, "hset", user2AdminHashKey, uid, adminId)
	return
}

// RemoveUserAdmin RemoveAdmin 移除用户客服
func (s *sChatRelation) RemoveUserAdmin(ctx gctx.Ctx, uid uint) (err error) {
	_, err = g.Redis().Do(ctx, "hdel", user2AdminHashKey, uid)
	return err
}

// GetUserValidAdmin GetValidAdmin 获取用户客服
func (s *sChatRelation) GetUserValidAdmin(ctx gctx.Ctx, uid uint) (uint, error) {
	val, err := g.Redis().HGet(ctx, user2AdminHashKey, strconv.Itoa(int(uid)))
	if err != nil {
		return 0, err
	}
	if val.IsNil() {
		return 0, nil
	}
	adminId := val.Uint()
	limitTime := s.GetLimitTime(ctx, adminId, uid)
	if limitTime > time.Now().Unix() {
		return val.Uint(), nil
	} else {
		// 无效了直接清除掉
		_ = s.RemoveUserAdmin(ctx, uid)
		return 0, nil
	}

}
