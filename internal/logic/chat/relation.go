package chat

import (
	"fmt"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"strconv"
	"time"
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
	relation                 = sChatRelation{}
)

type sChatRelation struct {
}

func (s *sChatRelation) getUserCacheKey(adminId uint) string {
	return fmt.Sprintf(adminChatUserKey, adminId)
}

// addUser 接入user
func (s *sChatRelation) addUser(ctx gctx.Ctx, adminId uint, uid uint) (err error) {
	err = s.setUserAdmin(ctx, uid, adminId)
	if err != nil {
		return
	}
	_, err = g.Redis().ZAdd(ctx, s.getUserCacheKey(adminId), nil, gredis.ZAddMember{
		Member: uid,
		Score:  float64(time.Now().Unix() + DefaultSessionTime),
	})
	if err != nil {
		return
	}
	err = s.updateUser(ctx, adminId, uid)
	return err
}

// updateUser 更新user
// 更新limit time
// 更新最后聊天时间
func (s *sChatRelation) updateUser(ctx gctx.Ctx, adminId uint, uid uint) (err error) {
	err = s.updateLimitTime(ctx, adminId, uid, DefaultSessionTime)
	if err != nil {
		return err
	}
	err = s.updateLastChatTime(ctx, adminId, uid)
	return err
}

// removeUser 移除user
func (s *sChatRelation) removeUser(ctx gctx.Ctx, adminId uint, uid uint) (err error) {
	err = s.removeUserAdmin(ctx, uid)
	if err != nil {
		return err
	}
	err = s.removeLastChatTime(ctx, adminId, uid)
	if err != nil {
		return err
	}
	_, err = g.Redis().ZRem(ctx, s.getUserCacheKey(adminId), uid)
	return err
}

// isUserValid 检查用户对于客服是否合法
func (s *sChatRelation) isUserValid(ctx gctx.Ctx, adminId uint, uid uint) (bool, error) {
	b, err := s.getLimitTime(ctx, adminId, uid)
	if err != nil {
		return false, nil
	}
	return b > time.Now().Unix(), nil
}

// user是否存在
func (s *sChatRelation) isUserExist(ctx gctx.Ctx, adminId uint, uid uint) (bool, error) {
	val, err := g.Redis().ZScore(ctx, s.getUserCacheKey(adminId), uid)
	if err != nil {
		return false, err
	}
	if val <= 0 {
		return false, nil
	}
	return true, nil
}

// 获取最后聊天时间
func (s *sChatRelation) getLastChatTime(ctx gctx.Ctx, adminId uint, uid uint) (uint, error) {
	val, err := g.Redis().HGet(ctx, fmt.Sprintf(adminUserLastChatKey, adminId), gconv.String(uid))
	if err != nil {
		return 0, err
	}
	return val.Uint(), nil
}

// removeLastChatTime 移除最后聊天时间
func (s *sChatRelation) removeLastChatTime(ctx gctx.Ctx, adminId uint, uid uint) error {
	_, err := g.Redis().HDel(ctx, fmt.Sprintf(adminUserLastChatKey, adminId), gconv.String(uid))
	return err
}

// updateLastChatTime 更新最后聊天时间
func (s *sChatRelation) updateLastChatTime(ctx gctx.Ctx, adminId uint, uid uint) error {
	_, err := g.Redis().HSet(ctx, fmt.Sprintf(adminUserLastChatKey, adminId), g.Map{
		gconv.String(uid): time.Now().Unix(),
	})
	return err
}

// getActiveCount 获取有效的用户数量
func (s *sChatRelation) getActiveCount(ctx gctx.Ctx, adminId uint) (uint, error) {
	val, err := g.Redis().ZRange(ctx,
		s.getUserCacheKey(adminId), time.Now().Unix(), -1, gredis.ZRangeOption{ByScore: true})
	if err != nil {
		return 0, err
	}
	return uint(len(val.Int64s())), nil
}

// updateLimitTime 更新有效期
func (s *sChatRelation) updateLimitTime(ctx gctx.Ctx, adminId uint, uid uint, duration int64) error {
	exist, err := s.isUserExist(ctx, adminId, uid)
	if err != nil {
		return err
	}
	if !exist {
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, "user not valid")
	}
	_, err = g.Redis().ZAdd(ctx, s.getUserCacheKey(adminId), nil, gredis.ZAddMember{
		Score:  float64(time.Now().Unix() + duration),
		Member: uid,
	})
	return err
}

// getLimitTime 获取有效期
func (s *sChatRelation) getLimitTime(ctx gctx.Ctx, adminId uint, uid uint) (int64, error) {
	val, err := g.Redis().ZScore(ctx, s.getUserCacheKey(adminId), uid)
	if err != nil {
		return 0, err
	}
	return gconv.Int64(val), nil
}

func (s *sChatRelation) getInvalidUsers(ctx gctx.Ctx, adminId uint) ([]uint, error) {
	val, err := g.Redis().ZRange(ctx,
		s.getUserCacheKey(adminId), 0, time.Now().Unix(), gredis.ZRangeOption{ByScore: true})
	if err != nil {
		return nil, err
	}
	return val.Uints(), nil
}

// getUsersWithLimitTime 获取所有user以及对应的有效期
func (s *sChatRelation) getUsersWithLimitTime(ctx gctx.Ctx, adminId uint) (uids []uint, times []int64, err error) {
	val, err := g.Redis().ZRevRange(ctx, s.getUserCacheKey(adminId), 0, -1, gredis.ZRevRangeOption{WithScores: true})
	if err != nil {
		return
	}
	for _, item := range val.Vars() {
		s := item.Slice()
		uids = append(uids, gconv.Uint(s[0]))
		times = append(times, gconv.Int64(s[1]))
	}
	return
}

// setUserAdmin SetAdmin 设置用户客服
func (s *sChatRelation) setUserAdmin(ctx gctx.Ctx, uid uint, adminId uint) (err error) {
	_, err = g.Redis().HSet(ctx, user2AdminHashKey, g.Map{gconv.String(uid): adminId})
	return
}

// removeUserAdmin RemoveAdmin 移除用户客服
func (s *sChatRelation) removeUserAdmin(ctx gctx.Ctx, uid uint) (err error) {
	_, err = g.Redis().HDel(ctx, user2AdminHashKey, gconv.String(uid))
	return err
}

// getUserValidAdmin GetValidAdmin 获取用户客服
func (s *sChatRelation) getUserValidAdmin(ctx gctx.Ctx, uid uint) (uint, error) {
	val, err := g.Redis().HGet(ctx, user2AdminHashKey, strconv.Itoa(int(uid)))
	if err != nil {
		return 0, err
	}
	if val.IsNil() {
		return 0, nil
	}
	adminId := val.Uint()
	limitTime, err := s.getLimitTime(ctx, adminId, uid)
	if err != nil {
		return 0, err
	}
	if limitTime > time.Now().Unix() {
		return val.Uint(), nil
	} else {
		// 无效了直接清除掉
		_ = s.removeUserAdmin(ctx, uid)
		return 0, nil
	}

}
