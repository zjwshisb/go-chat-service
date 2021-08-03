package chat

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"ws/app/databases"
	"ws/configs"
)

const (
	UserSubscribeKey = "user:%d:subscribe:%s"
)

// 标记 用户微信订阅消息 已订阅
func SetSubscribe(uid int64) error {
	ctx := context.Background()
	templateId := configs.Wechat.SubscribeTemplateIdOne
	key := fmt.Sprintf(UserSubscribeKey, uid, templateId)
	cmd := databases.Redis.Set(ctx, key, 1, 0)
	return cmd.Err()
}
// 查询 用户微信订阅消息
func IsSubScribe(uid int64) bool  {
	ctx := context.Background()
	templateId := configs.Wechat.SubscribeTemplateIdOne
	key := fmt.Sprintf(UserSubscribeKey, uid, templateId)
	cmd := databases.Redis.Get(ctx, key)
	if cmd.Err() == redis.Nil {
		return false
	}
	return true
}
// 删除 用户微信订阅消息 标记
func DelSubScribe(uid int64) bool {
	ctx := context.Background()
	templateId := configs.Wechat.SubscribeTemplateIdOne
	key := fmt.Sprintf(UserSubscribeKey, uid, templateId)
	databases.Redis.Del(ctx, key)
	return true
}


