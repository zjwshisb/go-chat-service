package backend

import (
	"context"
	api "gf-chat/api/v1/backend/admin"
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

func (cAdmin cAdmin) Index(ctx context.Context, req *api.ListReq) (res *api.ListRes, err error) {
	p, err := service.Admin().Paginate(ctx, &do.CustomerAdmins{
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
	}, model.QueryInput{
		Page:      req.Current,
		Size:      req.PageSize,
	}, nil, nil)
	if err != nil {
		return
	}
	item := slice.Map(p.Items, func(index int, item model.CustomerAdmin) api.ListItem {
		online := service.Chat().IsOnline(item.CustomerId, item.Id, "admin")
		var lastOnline *gtime.Time
		if item.Setting != nil && item.Setting.LastOnline.Year() != 1 {
			lastOnline = item.Setting.LastOnline
		}
		return api.ListItem{
			Id:            item.Id,
			Username:      item.Username,
			Online:        online,
			AcceptedCount: service.ChatRelation().GetActiveCount(item.Id),
			LastOnline:    lastOnline,
		}
	})
	res = &api.ListRes{
		Items: item,
		Total: p.Total,
	}
	return
}

func (cAdmin cAdmin) Show(ctx context.Context, req *api.DetailReq) (res *api.DetailRes, err error) {
	id := g.RequestFromCtx(ctx).GetRouter("id").Val()
	chart, model, err := service.Admin().GetDetail(ctx, id, req.Month)
	if err != nil {
		return
	}
	return &api.DetailRes{
		Chart: chart,
		Admin: model,
	}, nil
}
