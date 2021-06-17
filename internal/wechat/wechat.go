package wechat

import (
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram"
	"github.com/silenceper/wechat/v2/miniprogram/config"
	systemConfig "ws/configs"
)
var (
	mp *miniprogram.MiniProgram
)

func GetMp() *miniprogram.MiniProgram {
	if mp == nil {
		wc := wechat.NewWechat()
		memory := cache.NewMemory()
		cfg := &config.Config{
			AppID:     systemConfig.Wechat.MiniProgramAppId,
			AppSecret: systemConfig.Wechat.MiniProgramAppSecret,
			Cache: memory,
		}
		mp = wc.GetMiniProgram(cfg)
	}
	return mp
}
