package userCtx

import (
	"context"
	"gf-chat/internal/model"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/net/ghttp"
)

const userCtxKey = "user-ctx"

func init() {
	service.RegisterUserCtx(&sUserCtx{})
}

type sUserCtx struct {
}

// Init 初始化上下文对象指针到上下文对象中，以便后续的请求流程中可以修改。
func (s *sUserCtx) Init(r *ghttp.Request, customCtx *model.UserCtx) {
	r.SetCtxVar(userCtxKey, customCtx)
}

// Get 获得上下文变量，如果没有设置，那么返回nil
func (s *sUserCtx) Get(ctx context.Context) *model.UserCtx {
	value := ctx.Value(userCtxKey)
	if value == nil {
		return nil
	}
	if localCtx, ok := value.(*model.UserCtx); ok {
		return localCtx
	}
	return nil
}

// GetCustomerId 获取客户id
func (s *sUserCtx) GetCustomerId(ctx context.Context) int {
	admin := s.GetUser(ctx)
	if admin != nil {
		return admin.CustomerId
	}
	return 0
}

// GetUserApp 获取admin实体
func (s *sUserCtx) GetUserApp(ctx context.Context) *entity.UserApps {
	userCtx := s.Get(ctx)
	if userCtx != nil {
		return userCtx.UserApp
	}
	return nil
}

// GetUser 获取admin实体
func (s *sUserCtx) GetUser(ctx context.Context) *entity.Users {
	userCtx := s.Get(ctx)
	if userCtx != nil {
		return userCtx.Entity
	}
	return nil
}
