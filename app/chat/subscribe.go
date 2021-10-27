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

var SubScribeService = &subscribeService{}

type subscribeService struct {
}

func (subscribeService *subscribeService) Set(uid int64) error {
	ctx := context.Background()
	templateId := configs.Wechat.SubscribeTemplateIdOne
	key := fmt.Sprintf(UserSubscribeKey, uid, templateId)
	cmd := databases.Redis.Set(ctx, key, 1, 0)
	return cmd.Err()
}
func (subscribeService *subscribeService) IsSet(uid int64) bool {
	ctx := context.Background()
	templateId := configs.Wechat.SubscribeTemplateIdOne
	key := fmt.Sprintf(UserSubscribeKey, uid, templateId)
	cmd := databases.Redis.Get(ctx, key)
	if cmd.Err() == redis.Nil {
		return false
	}
	return true
}

func (subscribeService *subscribeService) Remove(uid int64)  bool {
	ctx := context.Background()
	templateId := configs.Wechat.SubscribeTemplateIdOne
	key := fmt.Sprintf(UserSubscribeKey, uid, templateId)
	databases.Redis.Del(ctx, key)
	return true
}


