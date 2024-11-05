package adminctx

import (
	"context"
	"gf-chat/internal/dao"
	"gf-chat/internal/model"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/model/relation"
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
func (s *sAdminCtx) GetKitchens(ctx context.Context) []*relation.CustomerKitchen {
	admin := s.Get(ctx)
	if admin == nil {
		return nil
	}
	if admin.GetKitchens() == nil {
		kitchens := make([]*relation.CustomerKitchen, 0)
		if admin.Entity.IsSuper == 1 {
			_ = dao.CustomerKitchens.Ctx(ctx).Where("customer_id", admin.Entity.CustomerId).WithAll().Scan(&kitchens)
		} else {
			kids, _ := g.Model("customer_admin_has_kitchen").
				Where("admin_id", admin.Entity.Id).
				Array("kitchen_id")
			_ = dao.CustomerKitchens.Ctx(ctx).WhereIn("id", kids).Scan(&kitchens)
			if len(kitchens) == 0 {
				_ = dao.CustomerKitchens.Ctx(ctx).Where("customer_id").WithAll().Scan(&kitchens)
			}
		}
		admin.SetKitchens(kitchens)
	}
	return admin.GetKitchens()
}

func (s *sAdminCtx) GetSchools(ctx context.Context) []*entity.Schools {
	admin := s.Get(ctx)
	if admin == nil {
		return nil
	}
	if admin.GetSchools() == nil {
		schools := make([]*entity.Schools, 0)
		for _, kitchen := range s.GetKitchens(ctx) {
			schools = append(schools, kitchen.Schools...)
		}
		unSetSchools := make([]*entity.Schools, 0)
		_ = dao.Schools.Ctx(ctx).Where("customer_id", admin.Entity.CustomerId).
			WhereNull("kitchen_id").Scan(&unSetSchools)
		schools = append(schools, unSetSchools...)
		admin.SetSchools(schools)
	}
	return admin.GetSchools()
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
