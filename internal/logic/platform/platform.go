package file

import (
	"context"
	"gf-chat/internal/consts"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/frame/g"
)

func init() {
	service.RegisterPlatform(&sPlatform{})
}

type sPlatform struct {
}

// GetPlatform 获取用户的平台
// todo
func (p sPlatform) GetPlatform(ctx context.Context) string {
	request := g.RequestFromCtx(ctx)
	_ = request.GetHeader("aaa")
	return consts.PlatformH5
}
