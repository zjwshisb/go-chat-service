package backend

import (
	"context"
	baseApi "gf-chat/api/v1"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/model"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/os/gtime"
)

var CCustomerAdmin = cCustomerAdmin{}

type cCustomerAdmin struct {
}

func (cAdmin cCustomerAdmin) Index(ctx context.Context, req *api.CustomerAdminListReq) (res *baseApi.ListRes[api.CustomerAdmin], err error) {
	w := g.Map{
		"customer_id": service.AdminCtx().GetCustomerId(ctx),
	}
	if req.Username != "" {
		w["username like"] = "%" + req.Username + "%"
	}
	p, err := service.Admin().Paginate(ctx, w, req.Paginate, g.Slice{model.CustomerAdmin{}.Setting}, nil)
	if err != nil {
		return
	}
	item := slice.Map(p.Items, func(_ int, item *model.CustomerAdmin) api.CustomerAdmin {
		online := service.Chat().IsOnline(item.CustomerId, item.Id, "admin")
		var lastOnline *gtime.Time
		if item.Setting != nil && item.Setting.LastOnline != nil {
			lastOnline = item.Setting.LastOnline
		}
		return api.CustomerAdmin{
			Id:            item.Id,
			Username:      item.Username,
			Online:        online,
			AcceptedCount: service.ChatRelation().GetActiveCount(ctx, item.Id),
			LastOnline:    lastOnline,
			UpdatedAt:     item.UpdatedAt,
			CreatedAt:     item.CreatedAt,
		}
	})
	res = baseApi.NewListResp(item, p.Total)
	return
}
