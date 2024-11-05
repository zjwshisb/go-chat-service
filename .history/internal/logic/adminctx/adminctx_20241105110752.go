package adminctx

import (
	"context"
	"gf-chat/internal/model"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

const adminCtxKey = "admin-ctx"

func init() {
	service.RegisterAdminCtx(New())
}

func New() *sAdminCtx {
	return &sAdminCtx{}
}

type sAdminCtx struct {
}

// Init 初始化上下文对象指针到上下文对象中，以便后续的请求流程中可以修改。
func (s *sAdminCtx) Init(r *ghttp.Request, customCtx *model.AdminCtx) {
	r.SetCtxVar(adminCtxKey, customCtx)
}

// Get 获得上下文变量，如果没有设置，那么返回nil
func (s *sAdminCtx) Get(ctx context.Context) *model.AdminCtx {
	value := ctx.Value(adminCtxKey)
	if value == nil {
		return nil
	}
	if localCtx, ok := value.(*model.AdminCtx); ok {
		return localCtx
	}
	return nil
}

// GetCustomerId 获取客户id
func (s *sAdminCtx) GetCustomerId(ctx context.Context) int {
	admin := s.GetAdmin(ctx)
	if admin != nil {
		return admin.CustomerId
	}
	return 0
}

// GetAdmin 获取admin实体
func (s *sAdminCtx) GetAdmin(ctx context.Context) *entity.CustomerAdmins {
	adminCtx := s.Get(ctx)
	if adminCtx != nil {
		return adminCtx.Entity
	}
	return nil
}

// SetUser 将上下文信息设置到上下文请求中
func (s *sAdminCtx) SetUser(ctx context.Context, Admin *entity.CustomerAdmins) {
	s.Get(ctx).Entity = Admin
}

// SetData 将上下文信息设置到上下文请求中
func (s *sAdminCtx) SetData(ctx context.Context, data g.Map) {
	s.Get(ctx).Data = data
}
