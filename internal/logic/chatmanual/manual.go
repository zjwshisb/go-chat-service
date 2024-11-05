package chat

import (
	"context"
	"fmt"
	"gf-chat/internal/service"
	"github.com/go-redis/redis/v8"
	"github.com/gogf/gf/v2/frame/g"
	"time"
)

func init() {
	service.RegisterChatManual(&sChatManual{})
}

const (
	manualUserKey = "user:%d:manual"
)

type sChatManual struct {
}

func (s *sChatManual) getManualKey(gid int) string {
	return fmt.Sprintf(manualUserKey, gid)
}

// Add 加入到待人工接入sortSet
func (s *sChatManual) Add(uid int, gid int) error {
	ctx := context.Background()
	_, err := g.Redis().Do(ctx, "zadd", s.getManualKey(gid), time.Now().Unix(), uid)
	return err
}

// IsIn 是否在待人工接入列表中
func (s *sChatManual) IsIn(uid int, customerId int) bool {
	ctx := context.Background()
	val, err := g.Redis().Do(ctx, "zrank", s.getManualKey(customerId), uid)
	// key在sort set 中不存在
	if err == redis.Nil {
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
func (s *sChatManual) Remove(uid int, customerId int) error {
	ctx := context.Background()
	_, err := g.Redis().Do(ctx, "ZRem", s.getManualKey(customerId), uid)
	return err
}

// GetTotalCount 获取待人工接入的数量
func (s *sChatManual) GetTotalCount(customerId int) int {
	ctx := context.Background()
	val, _ := g.Redis().Do(ctx, "ZCard", s.getManualKey(customerId))
	return val.Int()
}

// GetCountByTime 获取指定时间的数量
func (s *sChatManual) GetCountByTime(customerId int, min string, max string) int64 {
	ctx := context.Background()
	val, _ := g.Redis().Do(ctx, "ZCount", s.getManualKey(customerId), min, max)
	return val.Int64()
}

// GetByTime 通过加入时间获取
func (s *sChatManual) GetByTime(customerId int, min string, max string) []string {
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
func (s *sChatManual) GetTime(uid int, customerId int) float64 {
	ctx := context.Background()
	val, _ := g.Redis().Do(ctx, "ZScore", s.getManualKey(customerId), uid)
	return val.Float64()
}

// GetAll 获取所有待人工接入ids
func (s *sChatManual) GetAll(customerId int) []int64 {
	ctx := context.Background()
	val, _ := g.Redis().Do(ctx, "ZRangeByScore", s.getManualKey(customerId), &redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+info",
		Offset: 0,
		Count:  0,
	})
	return val.Int64s()
}

func (s *sChatManual) GetBySource(customerId int, Offset, count int) []int64 {
	ctx := context.Background()
	val, _ := g.Redis().Do(ctx, "ZRangeByScore", s.getManualKey(customerId), Offset, count)
	return val.Int64s()
}
