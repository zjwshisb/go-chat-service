package chat

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
)

// 此包内未返回的错误使用该logger
var log *glog.Logger

func init() {

	configMap, err := g.Config().Get(gctx.New(), "websocket.logger", g.Map{})
	if err != nil {
		panic(err)
	}
	log = glog.New()
	if len(configMap.Map()) > 0 {
		err = log.SetConfigWithMap(configMap.Map())
		if err != nil {
			panic(err)
		}
	}
}
