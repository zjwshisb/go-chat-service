package backend

import (
	"context"
	baseApi "gf-chat/api"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/dao"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"
)

var CSystemRule = &cSystemRule{}

type cSystemRule struct {
}

func (c cSystemRule) Index(ctx context.Context, req *api.SystemRuleListReq) (res *baseApi.NormalRes[api.SystemRuleListRes], err error) {
	paginator, err := service.AutoRule().Paginate(ctx, &do.CustomerChatAutoRules{
		IsSystem:   1,
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
	}, model.QueryInput{
		Size: 100,
		Page: 1,
	}, nil, nil)
	if err != nil {
		return
	}
	rr := api.SystemRuleListRes{}
	for _, item := range paginator.Items {
		rr = append(rr, api.SystemRuleListItem{
			Name:      item.Name,
			MessageId: item.MessageId,
			Id:        item.Id,
		})
	}
	return baseApi.NewResp(rr), nil
}

func (c cSystemRule) Update(ctx context.Context, req *api.SystemRuleUpdateReq) (res *baseApi.NilRes, err error) {
	_, err = dao.CustomerChatAutoRules.Ctx(ctx).Where(do.CustomerChatAutoRules{
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
		IsSystem:   1,
	}).Data("message_id", 0).Update()
	if err != nil {
		return
	}
	for sId, newId := range req.Data {
		_, err = dao.CustomerChatAutoRules.Ctx(ctx).Where(do.CustomerChatAutoRules{
			CustomerId: service.AdminCtx().GetCustomerId(ctx),
			IsSystem:   1,
			Id:         sId,
		}).Data("message_id", newId).Update()
		if err != nil {
			return
		}
	}
	return baseApi.NewNilResp(), nil
}
