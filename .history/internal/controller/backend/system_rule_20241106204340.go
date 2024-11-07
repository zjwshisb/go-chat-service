package backend

import (
	"context"
	baseApi "gf-chat/api"
	api "gf-chat/api/v1/backend/systemrule"
	"gf-chat/internal/dao"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"
)

var CSystemRule = &cSystemRule{}

type cSystemRule struct {
}

func (c cSystemRule) Index(ctx context.Context, req *api.ListReq) (res *api.ListRes, err error) {
	rules, _ := service.AutoRule().Paginate(ctx, &do.CustomerChatAutoRules{
		IsSystem:   1,
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
	}, model.QueryInput{
		Size:      100,
		Page:      1,
		WithTotal: false,
	})
	rr := api.ListRes{}
	for _, item := range rules {
		rr = append(rr, api.ListItem{
			Name:      item.Name,
			MessageId: item.MessageId,
			Id:        item.Id,
		})
	}
	res = &rr
	return
}

func (c cSystemRule) Update(ctx context.Context, req *api.UpdateReq) (res *baseApi.NilRes, err error) {
	dao.CustomerChatAutoRules.Ctx(ctx).Where(do.CustomerChatAutoRules{
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
		IsSystem:   1,
	}).Data("message_id", 0).Update()
	for sId, newId := range req.Data {
		dao.CustomerChatAutoRules.Ctx(ctx).Where(do.CustomerChatAutoRules{
			CustomerId: service.AdminCtx().GetCustomerId(ctx),
			IsSystem:   1,
			Id:         sId,
		}).Data("message_id", newId).Update()
	}
	return &baseApi.NilRes{}, nil
}
