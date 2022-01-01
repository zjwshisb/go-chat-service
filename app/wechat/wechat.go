package wechat

import (
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram"
	"github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/spf13/viper"
)

var (
	mp *miniprogram.MiniProgram
)

func GetMp() *miniprogram.MiniProgram {
	if mp == nil {
		wc := wechat.NewWechat()
		memory := cache.NewMemory()
		cfg := &config.Config{
			AppID:     viper.GetString("Wechat.MiniProgramAppId"),
			AppSecret:  viper.GetString("Wechat.MiniProgramAppSecret"),
			Cache: memory,
		}
		mp = wc.GetMiniProgram(cfg)
	}
	return mp
}
