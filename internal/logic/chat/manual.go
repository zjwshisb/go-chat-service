package chat

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gogf/gf/v2/frame/g"
)

const (
	manualUserKey = "user:%d:manual"
)

var manual = &sManual{}

type sManual struct {
}

func (s sManual) manualKey(gid uint) string {
	return fmt.Sprintf(manualUserKey, gid)
}

// AddToManualSet 加入到待人工接入sortSet
func (s sManual) addToSet(ctx context.Context, uid uint, gid uint) error {
	_, err := g.Redis().Do(ctx, "zadd", s.manualKey(gid), time.Now().Unix(), uid)
	return err
}

// IsInManualSet 是否在待人工接入列表中
func (s sManual) isInSet(ctx context.Context, uid uint, customerId uint) bool {
	val, err := g.Redis().Do(ctx, "zrank", s.manualKey(customerId), uid)
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

// RemoveManual 从待人工接入列表中移除
func (s sManual) removeFromSet(ctx context.Context, uid uint, customerId uint) error {
	_, err := g.Redis().Do(ctx, "ZRem", s.manualKey(customerId), uid)
	return err
}

// GetTotalCount 获取待人工接入的数量
func (s sManual) getCount(ctx context.Context, customerId uint) (count uint, err error) {
	val, err := g.Redis().Do(ctx, "ZCard", s.manualKey(customerId))
	if err != nil {
		return
	}
	count = val.Uint()
	return
}

// GetCountByTime 获取指定时间的数量
func (s sManual) getCountByTime(ctx context.Context, customerId uint, min string, max string) (count uint, err error) {
	val, err := g.Redis().Do(ctx, "ZCount", s.manualKey(customerId), min, max)
	if err != nil {
		return
	}
	count = val.Uint()
	return
}

// GetManualByTime 通过加入时间获取
func (s sManual) getListByTime(ctx context.Context, customerId uint, min string, max string) (uids []uint, err error) {
	val, err := g.Redis().Do(ctx, "ZRangeByScore", s.manualKey(customerId), &redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: 0,
		Count:  0,
	})
	if err != nil {
		return
	}
	uids = val.Uints()
	return
}

// GetManualTime 获取加入时间
func (s sManual) getAddTime(ctx context.Context, uid uint, customerId uint) (time float64, err error) {
	val, err := g.Redis().Do(ctx, "ZScore", s.manualKey(customerId), uid)
	if err != nil {
		return
	}
	time = val.Float64()
	return
}

// GetAll 获取所有待人工接入ids
func (s sManual) getAllList(ctx context.Context, customerId uint) (uids []uint, err error) {
	val, err := g.Redis().Do(ctx, "ZRangeByScore", s.manualKey(customerId), &redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+info",
		Offset: 0,
		Count:  0,
	})
	if err != nil {
		return
	}
	uids = val.Uints()
	return
}

func (s sManual) getList(ctx context.Context, customerId uint, Offset, count uint) (uids []uint, err error) {
	val, err := g.Redis().Do(ctx, "ZRangeByScore", s.manualKey(customerId), Offset, count)
	if err != nil {
		return
	}
	uids = val.Uints()
	return
}
