package chat

import (
	"context"
	"fmt"
	"strconv"
	"ws/app/databases"
)

const (
	// 用户 => 客服 hashes
	user2AdminHashKey = "user-to-admin"
)

var UserService = &userService{}


type userService struct {
}


// 设置客服
func (userService *userService) SetAdmin(uid int64, adminId int64) error {
	ctx := context.Background()
	cmd := databases.Redis.HSet(ctx, user2AdminHashKey, uid, adminId)
	return cmd.Err()
}

func (userService *userService) GetReqId(uid int64) int64 {
	key := fmt.Sprintf("user:%d:req-id", uid)
	ctx := context.Background()
	cmd := databases.Redis.Incr(ctx, key)
	return cmd.Val()
}

// 移除当前客服
func (userService *userService) RemoveAdmin(uid int64) error  {
	ctx := context.Background()
	cmd := databases.Redis.HDel(ctx, user2AdminHashKey, strconv.FormatInt(uid, 10))
	return cmd.Err()
}

// 获取客服
func (userService *userService) GetValidAdmin(uid int64) int64 {
	ctx := context.Background()
	key := strconv.FormatInt(uid, 10)
	cmd := databases.Redis.HGet(ctx, user2AdminHashKey, key)
	if adminId, err := cmd.Int64(); err == nil {
		if AdminService.IsUserValid(uid, adminId) {
			return adminId
		}
	}
	return 0
}