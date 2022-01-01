package chat

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"ws/app/databases"
)

const (
	UserSubscribeKey = "user:%d:subscribe:%s"
)

var SubScribeService = &subscribeService{}

type subscribeService struct {
}

// Set 设置微信订阅消息id
func (subscribeService *subscribeService) Set(uid int64) error {
	ctx := context.Background()
	templateId := viper.GetString("Wechat.SubscribeTemplateIdOne")
	key := fmt.Sprintf(UserSubscribeKey, uid, templateId)
	cmd := databases.Redis.Set(ctx, key, 1, 0)
	return cmd.Err()
}

// IsSet 是否设置微信订阅消息id
func (subscribeService *subscribeService) IsSet(uid int64) bool {
	ctx := context.Background()
	templateId := viper.GetString("Wechat.SubscribeTemplateIdOne")
	key := fmt.Sprintf(UserSubscribeKey, uid, templateId)
	cmd := databases.Redis.Get(ctx, key)
	if cmd.Err() == redis.Nil {
		return false
	}
	return true
}

// Remove 移除微信订阅消息id
func (subscribeService *subscribeService) Remove(uid int64)  bool {
	ctx := context.Background()
	templateId := viper.GetString("Wechat.SubscribeTemplateIdOne")
	key := fmt.Sprintf(UserSubscribeKey, uid, templateId)
	databases.Redis.Del(ctx, key)
	return true
}


