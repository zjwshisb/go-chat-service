package backend

import (
	"context"
	baseApi "gf-chat/api"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/os/gtime"
)

var CCustomerAdmin = cCustomerAdmin{}

type cCustomerAdmin struct {
}

func (cAdmin cCustomerAdmin) Index(ctx context.Context, req *api.CustomerAdminListReq) (res *baseApi.ListRes[api.CustomerAdmin], err error) {
	p, err := service.Admin().Paginate(ctx, &do.CustomerAdmins{
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
	}, req.Paginate, nil, nil)
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

func (cAdmin cCustomerAdmin) Show(ctx context.Context, req *api.CustomerAdminDetailReq) (res *baseApi.NormalRes[api.AdminDetailRes], err error) {

	return
}