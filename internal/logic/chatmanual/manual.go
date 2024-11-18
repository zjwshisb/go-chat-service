package chat

import (
	"context"
	"errors"
	"fmt"
	"gf-chat/internal/service"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gogf/gf/v2/frame/g"
)

func init() {
	service.RegisterChatManual(&sChatManual{})
}

const (
	manualUserKey = "user:%d:manual"
)

type sChatManual struct {
}

func (s *sChatManual) getManualKey(gid uint) string {
	return fmt.Sprintf(manualUserKey, gid)
}

// Add 加入到待人工接入sortSet
func (s *sChatManual) Add(uid uint, gid uint) error {
	ctx := context.Background()
	_, err := g.Redis().Do(ctx, "zadd", s.getManualKey(gid), time.Now().Unix(), uid)
	return err
}

// IsIn 是否在待人工接入列表中
func (s *sChatManual) IsIn(uid uint, customerId uint) bool {
	ctx := context.Background()
	val, err := g.Redis().Do(ctx, "zrank", s.getManualKey(customerId), uid)
	// key在sort set 中不存在
	if errors.Is(err, redis.Nil) {
		return false
	}
	// sort set 的key 不存在
	if val.Val() == nil {
		return false
	}
	if val.Val() == 0 {
		return false
	}
	return true
}

// Remove 从待人工接入列表中移除
func (s *sChatManual) Remove(uid uint, customerId uint) error {
	ctx := context.Background()
	_, err := g.Redis().Do(ctx, "ZRem", s.getManualKey(customerId), uid)
	return err
}

// GetTotalCount 获取待人工接入的数量
func (s *sChatManual) GetTotalCount(customerId uint) uint {
	ctx := context.Background()
	val, _ := g.Redis().Do(ctx, "ZCard", s.getManualKey(customerId))
	return val.Uint()
}

// GetCountByTime 获取指定时间的数量
func (s *sChatManual) GetCountByTime(customerId uint, min string, max string) uint {
	ctx := context.Background()
	val, _ := g.Redis().Do(ctx, "ZCount", s.getManualKey(customerId), min, max)
	return val.Uint()
}

// GetByTime 通过加入时间获取
func (s *sChatManual) GetByTime(customerId uint, min string, max string) []string {
	ctx := context.Background()
	val, _ := g.Redis().Do(ctx, "ZRangeByScore", s.getManualKey(customerId), &redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: 0,
		Count:  0,
	})
	return val.Strings()
}

// GetTime 获取加入时间
func (s *sChatManual) GetTime(uid uint, customerId uint) float64 {
	ctx := context.Background()
	val, _ := g.Redis().Do(ctx, "ZScore", s.getManualKey(customerId), uid)
	return val.Float64()
}

// GetAll 获取所有待人工接入ids
func (s *sChatManual) GetAll(customerId uint) []uint {
	ctx := context.Background()
	val, _ := g.Redis().Do(ctx, "ZRangeByScore", s.getManualKey(customerId), &redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+info",
		Offset: 0,
		Count:  0,
	})
	return val.Uints()
}

func (s *sChatManual) GetBySource(customerId uint, Offset, count uint) []uint {
	ctx := context.Background()
	val, _ := g.Redis().Do(ctx, "ZRangeByScore", s.getManualKey(customerId), Offset, count)
	return val.Uints()
}
