package backend

import (
	"context"
	"gf-chat/api/v1/backend"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/relation"
	"gf-chat/internal/service"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"time"
)

var CAdmin = cAdmin{}

type cAdmin struct {
}

func (cAdmin cAdmin) Index(ctx context.Context, req *backend.AdminIndexReq) (res *backend.AdminIndexRes, err error) {
	admins, total := service.Admin().Paginate(ctx, &do.CustomerAdmins{
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
	}, model.QueryInput{
		Page:      req.Current,
		Size:      req.PageSize,
		WithTotal: true,
	})
	item := slice.Map(admins, func(index int, item *relation.CustomerAdmins) backend.AdminListItem {
		online := service.Chat().IsOnline(item.CustomerId, item.Id, "admin")
		var lastOnline *time.Time
		if item.Setting != nil && item.Setting.LastOnline.Year() != 1 {
			lastOnline = &item.Setting.LastOnline
		}
		return backend.AdminListItem{
			Id:            item.Id,
			Username:      item.Username,
			Online:        online,
			AcceptedCount: service.ChatRelation().GetActiveCount(gconv.Int(item.Id)),
			LastOnline:    lastOnline,
		}
	})
	res = &backend.AdminIndexRes{
		Items: item,
		Total: total,
	}
	return
}

func (cAdmin cAdmin) Show(ctx context.Context, req *backend.AdminDetailReq) (res *backend.AdminDetailRes, err error) {
	id := g.RequestFromCtx(ctx).GetRouter("id").Val()
	chart, admin, err := service.Admin().GetDetail(ctx, id, req.Month)
	if err != nil {
		return
	}
	return &backend.AdminDetailRes{
		Chart: chart,
		Admin: admin,
	}, nil
}
