package backend

import (
	"context"
	baseApi "gf-chat/api"
	adminApi "gf-chat/api/v1/backend"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/os/gtime"
)

var CAdmin = cAdmin{}

type cAdmin struct {
}

func (cAdmin cAdmin) Index(ctx context.Context, req *adminApi.AdminListReq) (res *baseApi.ListRes[adminApi.AdminListItem], err error) {
	p, err := service.Admin().Paginate(ctx, &do.CustomerAdmins{
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
	}, req.Paginate, nil, nil)
	if err != nil {
		return
	}
	item := slice.Map(p.Items, func(index int, item model.CustomerAdmin) adminApi.AdminListItem {
		online := service.Chat().IsOnline(item.CustomerId, item.Id, "admin")
		var lastOnline *gtime.Time
		if item.Setting != nil && item.Setting.LastOnline.Year() != 1 {
			lastOnline = item.Setting.LastOnline
		}
		return adminApi.AdminListItem{
			Id:            item.Id,
			Username:      item.Username,
			Online:        online,
			AcceptedCount: service.ChatRelation().GetActiveCount(ctx, item.Id),
			LastOnline:    lastOnline,
		}
	})
	res = baseApi.NewListResp(item, p.Total)
	return
}

func (cAdmin cAdmin) Show(ctx context.Context, req *adminApi.AdminDetailReq) (res *baseApi.NormalRes[adminApi.AdminDetailRes], err error) {

	return
}
