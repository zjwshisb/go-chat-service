package backend

import (
	"context"
	baseApi "gf-chat/api"
	chatApi "gf-chat/api/v1/backend"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"

	"github.com/duke-git/lancet/v2/slice"
)

var CTransfer = &cTransfer{}

type cTransfer struct {
}

func (c cTransfer) Cancel(ctx context.Context, req *chatApi.TransferCancelReq) (resp *baseApi.NilRes, err error) {
	admin := service.AdminCtx().GetAdmin(ctx)
	transfer, err := service.ChatTransfer().First(ctx, do.CustomerChatTransfers{
		CustomerId: admin.CustomerId,
		CanceledAt: nil,
		AcceptedAt: nil,
	})
	if err != nil {
		return
	}
	service.ChatTransfer().Cancel(transfer)
	return &baseApi.NilRes{}, nil
}

func (c cTransfer) Index(ctx context.Context, req *chatApi.TransferListReq) (res *baseApi.ListRes[chatApi.ChatTransfer], err error) {
	customerId := service.AdminCtx().GetCustomerId(ctx)
	w := do.CustomerChatTransfers{
		CustomerId: customerId,
	}
	if req.Username != "" {
		uW := make(map[string]any)
		uW["phone"] = req.Username
		uW["customer_id"] = customerId
		users := service.User().GetUsers(ctx, do.Users{
			Username:   req.Username,
			CustomerId: customerId,
		})
		uids := slice.Map(users, func(index int, item *entity.Users) uint {
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
	p, err := service.ChatTransfer().Paginate(ctx, &w, model.QueryInput{
		Size: req.PageSize,
		Page: req.Current,
	}, nil, nil)
	if err != nil {
		return
	}
	items := slice.Map(p.Items, func(index int, item model.CustomerChatTransfer) chatApi.ChatTransfer {
		return service.ChatTransfer().ToChatTransfer(&item)
	})
	return baseApi.NewListResp(items, p.Total), nil
}
