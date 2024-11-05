package chatnotice

import "gf-chat/internal/service"

const (
	UserSubscribeKey = "user:%d:subscribe:%s"
)

func init() {
	service.RegisterChatNotice(&sChatNotice{})
}

type sChatNotice struct {
}

// Set 设置微信订阅消息id
func (s *sChatNotice) Set(uid int) error {
	//ctx := context.Background()
	//templateId := viper.GetString("Wechat.SubscribeTemplateIdOne")
	//key := fmt.Sprintf(UserSubscribeKey, uid, templateId)
	//cmd := databases.Redis.Set(ctx, key, 1, 0)
	return nil
}

// IsSet 是否设置微信订阅消息id
func (s *sChatNotice) IsSet(uid int) bool {
	//ctx := context.Background()
	//templateId := viper.GetString("Wechat.SubscribeTemplateIdOne")
	//key := fmt.Sprintf(UserSubscribeKey, uid, templateId)
	//cmd := databases.Redis.Get(ctx, key)
	//if cmd.Err() == redis.Nil {
	//	return false
	//}
	return true
}

// Remove 移除微信订阅消息id
func (s *sChatNotice) Remove(uid int) bool {
	//ctx := context.Background()
	//templateId := viper.GetString("Wechat.SubscribeTemplateIdOne")
	//key := fmt.Sprintf(UserSubscribeKey, uid, templateId)
	//databases.Redis.Del(ctx, key)
	return true
}
