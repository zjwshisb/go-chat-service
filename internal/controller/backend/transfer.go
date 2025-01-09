package backend

import (
	"context"
	baseApi "gf-chat/api"
	"gf-chat/api/backend/v1"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/duke-git/lancet/v2/slice"
)

var CTransfer = &cTransfer{}

type cTransfer struct {
}

func (c cTransfer) Cancel(ctx context.Context, _ *v1.TransferCancelReq) (resp *baseApi.NilRes, err error) {
	transfer, err := service.ChatTransfer().First(ctx, do.CustomerChatTransfers{
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
		Id:         ghttp.RequestFromCtx(ctx).GetRouter("id").String(),
	})
	if err != nil {
		return
	}
	err = service.ChatTransfer().Cancel(ctx, transfer)
	if err != nil {
		return
	}
	return baseApi.NewNilResp(), nil
}

func (c cTransfer) Index(ctx context.Context, req *v1.TransferListReq) (res *baseApi.ListRes[v1.ChatTransfer], err error) {
	customerId := service.AdminCtx().GetCustomerId(ctx)
	w := do.CustomerChatTransfers{
		CustomerId: customerId,
	}
	if req.Username != "" {
		uW := make(map[string]any)
		uW["username"] = req.Username
		uW["customer_id"] = customerId
		users, err := service.User().All(ctx, do.Users{
			Username:   req.Username,
			CustomerId: customerId,
		}, nil, nil)
		if err != nil {
			return nil, err
		}
		uids := slice.Map(users, func(index int, item *model.User) uint {
			return item.Id
		})
		w.UserId = uids
	}
	if req.ToAdminName != "" {
		admins, err := service.Admin().All(ctx, do.CustomerAdmins{
			Username:   req.ToAdminName,
			CustomerId: customerId,
		}, nil, nil)
		if err != nil {
			return nil, err
		}
		adminIds := slice.Map(admins, func(index int, item *model.CustomerAdmin) uint {
			return item.Id
		})
		w.ToAdminId = adminIds
	}
	if req.FromAdminName != "" {
		admins, err := service.Admin().All(ctx, do.CustomerAdmins{
			Username:   req.FromAdminName,
			CustomerId: customerId,
		}, nil, nil)
		if err != nil {
			return nil, err
		}
		adminIds := slice.Map(admins, func(index int, item *model.CustomerAdmin) uint {
			return item.Id
		})
		w.FromAdminId = adminIds
	}
	p, err := service.ChatTransfer().Paginate(ctx, &w, req.Paginate, g.Slice{
		model.CustomerChatTransfer{}.FormAdmin,
		model.CustomerChatTransfer{}.ToAdmin,
		model.CustomerChatTransfer{}.User,
	}, nil)
	if err != nil {
		return
	}
	items := slice.Map(p.Items, func(index int, item *model.CustomerChatTransfer) v1.ChatTransfer {
		return service.ChatTransfer().ToApi(item)
	})
	return baseApi.NewListResp(items, p.Total), nil
}
