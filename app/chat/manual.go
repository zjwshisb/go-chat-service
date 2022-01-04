package chat

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
	"ws/app/databases"
)

const (
	manualUserKey = "user:%d:manual"
)
var (
	ManualService = &manualService{}
)
type manualService struct {
}

func (manual manualService) getManualKey(gid int64) string {
	return fmt.Sprintf(manualUserKey, gid)
}

// Add 加入到待人工接入sortSet
func (manual *manualService) Add(uid int64, gid int64) error  {
	ctx := context.Background()
	cmd := databases.Redis.ZAdd(ctx, manual.getManualKey(gid), &redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: uid,
	})
	return cmd.Err()
}

// IsIn 是否在待人工接入列表中
func (manual *manualService) IsIn(uid int64, gid int64) bool {
	ctx := context.Background()
	cmd := databases.Redis.ZRank(ctx, manual.getManualKey(gid), strconv.FormatInt(uid, 10))
	if cmd.Err() == redis.Nil {
		return false
	}
	return true
}

// Remove 从待人工接入列表中移除
func (manual *manualService) Remove(uid int64, gid int64) error {
	ctx := context.Background()
	cmd := databases.Redis.ZRem(ctx, manual.getManualKey(gid), uid)
	return cmd.Err()
}

// GetTotalCount 获取待人工接入的数量
func (manual *manualService) GetTotalCount(gid int64) int64 {
	ctx := context.Background()
	cmd := databases.Redis.ZCard(ctx, manual.getManualKey(gid))
	return cmd.Val()
}

// GetCountByTime 获取指定时间的数量
func (manual *manualService) GetCountByTime(gid int64, min string, max string)  int64 {
	ctx := context.Background()
	cmd := databases.Redis.ZCount(ctx,manual.getManualKey(gid), min, max)
	return cmd.Val()
}

// GetByTime 通过加入时间获取
func (manual *manualService) GetByTime(gid int64, min string, max string) []string {
	ctx := context.Background()
	cmd := databases.Redis.ZRangeByScore(ctx, manual.getManualKey(gid), &redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: 0,
		Count:  0,
	})
	return cmd.Val()
}

// GetTime 获取加入时间
func (manual *manualService) GetTime(uid int64, gid int64) float64 {
	ctx := context.Background()
	cmd := databases.Redis.ZScore(ctx, manual.getManualKey(gid), strconv.FormatInt(uid, 10))
	return cmd.Val()
}

// GetAll 获取所有待人工接入ids
func (manual *manualService) GetAll(gid int64) []int64 {
	ctx := context.Background()
	cmd := databases.Redis.ZRangeByScore(ctx, manual.getManualKey(gid), &redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+info",
		Offset: 0,
		Count:  0,
	})
	uid := make([]int64, 0, len(cmd.Val()))
	for _, uidStr := range cmd.Val() {
		id , err := strconv.ParseInt(uidStr, 10, 64)
		if err == nil {
			uid = append(uid, id)
		}
	}
	return uid
}

