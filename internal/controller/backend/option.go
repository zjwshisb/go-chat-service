package backend

import (
	"context"
	baseApi "gf-chat/api"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/consts"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"
	"github.com/duke-git/lancet/v2/slice"
)

var COption = &cOption{}

type cOption struct {
}

func (c *cAutoMessage) AutoMessage(ctx context.Context, req *api.OptionAutoMessageReq) (res *baseApi.OptionRes, err error) {
	items, err := service.AutoMessage().All(ctx, do.CustomerChatAutoMessages{
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
	}, nil, nil)
	if err != nil {
		return
	}
	options := slice.Map(items, func(index int, item *model.CustomerChatAutoMessage) baseApi.Option {
		return baseApi.Option{
			Label: item.Name,
			Value: item.Id,
		}
	})
	return baseApi.NewOptionResp(options), nil
}

func (c *cAutoMessage) AutoRuleScene(ctx context.Context, req *api.OptionAutoRuleSceneReq) (res *baseApi.OptionRes, err error) {
	options := []baseApi.Option{
		{
			Label: "人工未接入",
			Value: consts.AutoRuleSceneNotAccepted,
		},
		{
			Label: "已接入但客服离线",
			Value: consts.AutoRuleSceneAdminOnline,
		},
		{
			Label: "已接入客服在线",
			Value: consts.AutoRuleSceneAdminOnline,
		},
	}
	return baseApi.NewOptionResp(options), nil
}
