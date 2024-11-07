package backend

import (
	"context"
	baseApi "gf-chat/api"
	api "gf-chat/api/v1/backend/transfer"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"

	"github.com/duke-git/lancet/v2/slice"
)

var CTransfer = &cTransfer{}

type cTransfer struct {
}

func (c cTransfer) Cancel(ctx context.Context, req *api.CancelReq) (*baseApi.NilRes, error) {
	admin := service.AdminCtx().GetAdmin(ctx)
	transfer, _ := service.ChatTransfer().First(do.CustomerChatTransfers{
		CustomerId: admin.CustomerId,
		CanceledAt: nil,
		AcceptedAt: nil,
	})
	if transfer != nil {
		service.ChatTransfer().Cancel(transfer)
	}
	return &baseApi.NilRes{}, nil
}

func (c cTransfer) Index(ctx context.Context, req *api.ListReq) (res *api.ListRes, err error) {
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
		})
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
		})
		if err != nil {
			return nil, err
		}
		adminIds := slice.Map(admins, func(index int, item *model.CustomerAdmin) uint {
			return item.Id
		})
		w.FromAdminId = adminIds
	}
	transfers, total := service.ChatTransfer().Paginate(ctx, &w, model.QueryInput{
		Size:      req.PageSize,
		Page:      req.Current,
		WithTotal: true,
	})
	items := slice.Map(transfers, func(index int, item *model.CustomerChatTransfer) model.ChatTransfer {
		return service.ChatTransfer().ToChatTransfer(item)
	})
	return &api.ListRes{
		Total: total,
		Items: items,
	}, nil

}
