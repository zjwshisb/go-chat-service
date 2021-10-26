package chat

import (
	"context"
	"strconv"
	"ws/app/databases"
)

var (
	ManualService = &manualService{}
)
type manualService struct {

}
func (manual *manualService) Add(uid int64) error  {
	ctx := context.Background()
	cmd := databases.Redis.SAdd(ctx, manualUserKey, uid)
	return cmd.Err()
}
func (manual *manualService) IsIn(uid int64) bool {
	ctx := context.Background()
	cmd := databases.Redis.SIsMember(ctx, manualUserKey, uid)
	return cmd.Val()
}
func (manual *manualService) Remove(uid int64) error {
	ctx := context.Background()
	cmd := databases.Redis.SRem(ctx, manualUserKey, uid)
	return cmd.Err()
}
func (manual *manualService) GetAll() []int64 {
	ctx := context.Background()
	cmd := databases.Redis.SMembers(ctx, manualUserKey)
	uid := make([]int64, 0, len(cmd.Val()))
	for _, uidStr := range cmd.Val() {
		id , err := strconv.ParseInt(uidStr, 10, 64)
		if err == nil {
			uid = append(uid, id)
		}
	}
	return uid
}

