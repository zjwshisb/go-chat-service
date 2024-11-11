package backend

import (
	"context"
	baseApi "gf-chat/api"
	adminApi "gf-chat/api/v1/backend/admin"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

var CAdmin = cAdmin{}

type cAdmin struct {
}

func (cAdmin cAdmin) Index(ctx context.Context, req *adminApi.ListReq) (res *baseApi.ListRes[adminApi.ListItem], err error) {
	p, err := service.Admin().Paginate(ctx, &do.CustomerAdmins{
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
	}, model.QueryInput{
		Page: req.Current,
		Size: req.PageSize,
	}, nil, nil)
	if err != nil {
		return
	}
	item := slice.Map(p.Items, func(index int, item model.CustomerAdmin) adminApi.ListItem {
		online := service.Chat().IsOnline(item.CustomerId, item.Id, "admin")
		var lastOnline *gtime.Time
		if item.Setting != nil && item.Setting.LastOnline.Year() != 1 {
			lastOnline = item.Setting.LastOnline
		}
		return adminApi.ListItem{
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

func (cAdmin cAdmin) Show(ctx context.Context, req *adminApi.DetailReq) (res *baseApi.NormalRes[adminApi.DetailRes], err error) {
	id := g.RequestFromCtx(ctx).GetRouter("id").Val()
	chart, model, err := service.Admin().GetDetail(ctx, id, req.Month)
	if err != nil {
		return
	}
	res = baseApi.NewResp(adminApi.DetailRes{
		Chart: chart,
		Admin: model,
	})
	return
}
