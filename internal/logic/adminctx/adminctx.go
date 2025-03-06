package adminctx

import (
	"context"
	"gf-chat/internal/model"
	"gf-chat/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

const adminCtxKey = "admin-ctx"

func init() {
	service.RegisterAdminCtx(newI())
}

func newI() *sAdminCtx {
	return &sAdminCtx{}
}

type sAdminCtx struct {
}

// Init initializes the admin context for a request
// Sets the admin context in the request context
func (s *sAdminCtx) Init(r *ghttp.Request, customCtx *model.AdminCtx) {
	r.SetCtxVar(adminCtxKey, customCtx)
}

// Get retrieves the admin context from the request context
// Returns nil if the context is not set
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

// GetCustomerId retrieves the customer ID from the admin context
// Returns 0 if the admin context is not set or the customer ID is not set
func (s *sAdminCtx) GetCustomerId(ctx context.Context) uint {
	admin := s.GetUser(ctx)
	if admin != nil {
		return admin.CustomerId
	}
	return 0
}

// GetId retrieves the admin entity ID from the admin context
// Returns 0 if the admin context is not set or the admin entity ID is not set
func (s *sAdminCtx) GetId(ctx context.Context) uint {
	adminCtx := s.Get(ctx)
	if adminCtx != nil {
		return adminCtx.Entity.Id
	}
	return 0
}

// GetUser retrieves the admin entity from the admin context
// Returns nil if the admin context is not set or the admin entity is not set
func (s *sAdminCtx) GetUser(ctx context.Context) *model.CustomerAdmin {
	adminCtx := s.Get(ctx)
	if adminCtx != nil {
		return adminCtx.Entity
	}
	return nil
}

// SetUser sets the admin entity in the admin context
func (s *sAdminCtx) SetUser(ctx context.Context, Admin *model.CustomerAdmin) {
	s.Get(ctx).Entity = Admin
}

// SetData sets the data in the admin context
func (s *sAdminCtx) SetData(ctx context.Context, data g.Map) {
	s.Get(ctx).Data = data
}
