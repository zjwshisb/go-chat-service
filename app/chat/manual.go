package chat

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
	"ws/app/databases"
)

const (
	manualUserKey = "user:manual"
)
var (
	ManualService = &manualService{}
)
type manualService struct {
}
// 加入到待人工接入sortSet
func (manual *manualService) Add(uid int64) error  {
	ctx := context.Background()
	cmd := databases.Redis.ZAdd(ctx, manualUserKey, &redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: uid,
	})
	return cmd.Err()
}
// 是否在待人工接入列表中
func (manual *manualService) IsIn(uid int64) bool {
	ctx := context.Background()
	cmd := databases.Redis.ZRank(ctx, manualUserKey, strconv.FormatInt(uid, 10))
	if cmd.Err() == redis.Nil {
		return false
	}
	return true
}
// 从待人工接入列表中移除
func (manual *manualService) Remove(uid int64) error {
	ctx := context.Background()
	cmd := databases.Redis.ZRem(ctx, manualUserKey, uid)
	return cmd.Err()
}
// 获取待人工接入的数量
func (manual *manualService) GetTotalCount() int64 {
	ctx := context.Background()
	cmd := databases.Redis.ZCard(ctx, manualUserKey)
	return cmd.Val()
}
// 通过加入时间获取
func (manual *manualService) getByTime(min string, max string) []string {
	ctx := context.Background()
	cmd := databases.Redis.ZRangeByScore(ctx, manualUserKey, &redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: 0,
		Count:  0,
	})
	return cmd.Val()
}
// 获取加入时间
func (manual *manualService) getTime(uid int64) float64 {
	ctx := context.Background()
	cmd := databases.Redis.ZScore(ctx, manualUserKey, strconv.FormatInt(uid, 10))
	return cmd.Val()
}
// 获取所有待人工接入ids
func (manual *manualService) GetAll() []int64 {
	ctx := context.Background()
	cmd := databases.Redis.ZRangeByScore(ctx, manualUserKey, &redis.ZRangeBy{
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

