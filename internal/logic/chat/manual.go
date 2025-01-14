package chat

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/util/gconv"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

const (
	manualUserKey = "user:%d:manual"
)

var manual = &iManual{}

type iManual struct {
}

func (s iManual) manualKey(customerId uint) string {
	return fmt.Sprintf(manualUserKey, customerId)
}

// AddToManualSet 加入到待人工接入sortSet
func (s iManual) addToSet(ctx context.Context, uid uint, customerId uint) error {
	_, err := g.Redis().ZAdd(ctx, s.manualKey(customerId), nil, gredis.ZAddMember{
		Score:  float64(time.Now().Unix()),
		Member: gconv.String(uid),
	})
	return err
}

// 是否在待人工接入列表中
func (s iManual) isInSet(ctx context.Context, uid uint, customerId uint) (bool, error) {
	val, err := g.Redis().ZScore(ctx, s.manualKey(customerId), uid)
	if err != nil {
		return false, err
	}
	if val == 0 {
		return false, nil
	}

	return true, nil
}

// RemoveManual 从待人工接入列表中移除
func (s iManual) removeFromSet(ctx context.Context, uid uint, customerId uint) error {
	_, err := g.Redis().ZRem(ctx, s.manualKey(customerId), uid)
	return err
}

// GetTotalCount 获取待人工接入的数量
func (s iManual) getCount(ctx context.Context, customerId uint) (count uint, err error) {
	val, err := g.Redis().ZCard(ctx, s.manualKey(customerId))
	if err != nil {
		return
	}
	return gconv.Uint(val), nil
}

// GetCountByTime 获取指定时间的数量
func (s iManual) getCountByTime(ctx context.Context, customerId uint, min string, max string) (count uint, err error) {
	val, err := g.Redis().ZCount(ctx, s.manualKey(customerId), min, max)
	if err != nil {
		return
	}
	return gconv.Uint(val), nil
}

// GetManualTime 获取加入时间
func (s iManual) getAddTime(ctx context.Context, uid uint, customerId uint) (time float64, err error) {
	time, err = g.Redis().ZScore(ctx, s.manualKey(customerId), uid)
	return
}

// GetAll 获取所有待人工接入ids
func (s iManual) getAllList(ctx context.Context, customerId uint) (uids []uint, err error) {
	val, err := g.Redis().ZRange(ctx, s.manualKey(customerId), 0, -1)
	if err != nil {
		return
	}
	uids = val.Uints()
	return
}

func (s iManual) getList(ctx context.Context, customerId uint, Offset, count uint) (uids []uint, err error) {
	val, err := g.Redis().ZRange(ctx, s.manualKey(customerId), int64(Offset), int64(count), gredis.ZRangeOption{WithScores: true})
	if err != nil {
		return
	}
	uids = val.Uints()
	return
}
