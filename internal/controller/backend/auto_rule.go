package backend

import (
	"context"
	baseApi "gf-chat/api/v1"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

var CAutoRule = &cAutoRule{}

type cAutoRule struct {
}

func (c cAutoRule) Index(ctx context.Context, req *api.AutoRuleListReq) (res *baseApi.ListRes[api.AutoRule], err error) {
	w := g.Map{
		"customer_id": service.AdminCtx().GetCustomerId(ctx),
		"is_system":   0,
	}
	name := strutil.Trim(req.Name)
	if name != "" {
		w["name like ?"] = "%" + name + "%"
	}
	if req.ReplyType != "" {
		w["reply_type"] = req.ReplyType
	}
	if req.ReplyType != "" {
		w["match_type"] = req.MatchType
	}
	if req.IsOpen != nil {
		w["is_open"] = *req.IsOpen
	}
	paginator, err := service.AutoRule().Paginate(ctx, w, req.Paginate, nil, nil)
	if err != nil {
		return
	}
	messageIds := slice.Map(paginator.Items, func(index int, item *model.CustomerChatAutoRule) uint {
		return item.MessageId
	})

	items := make([]api.AutoRule, len(paginator.Items))
	messages, err := service.AutoMessage().All(ctx, do.CustomerChatAutoMessages{
		Id: messageIds,
	}, nil, nil)
	if err != nil {
		return
	}
	apiMessages, err := service.AutoMessage().ToApis(ctx, messages)
	apiMessagesMap := slice.KeyBy(apiMessages, func(item *api.AutoMessage) uint {
		return item.Id
	})
	for index, i := range paginator.Items {
		apiRule := api.AutoRule{
			Id:        i.Id,
			Name:      i.Name,
			Match:     i.Match,
			MatchType: i.MatchType,
			ReplyType: i.ReplyType,
			IsOpen:    i.IsOpen,
			Scenes:    i.Scenes,
			Sort:      i.Sort,
			Count:     uint(i.Count),
			CreatedAt: i.CreatedAt,
			UpdatedAt: i.UpdatedAt,
		}
		if i.MessageId != 0 {
			msg := maputil.GetOrDefault(apiMessagesMap, i.MessageId, nil)
			apiRule.Message = msg
		}
		items[index] = apiRule

	}
	return baseApi.NewListResp(items, paginator.Total), nil
}
func (c cAutoRule) Form(ctx context.Context, _ *api.AutoRuleFormReq) (res *baseApi.NormalRes[api.AutoRuleForm], err error) {
	rule, err := service.AutoRule().First(ctx, do.CustomerChatAutoRules{
		IsSystem:   0,
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
		Id:         ghttp.RequestFromCtx(ctx).GetRouter("id").Val(),
	})
	if err != nil {
		return
	}
	res = baseApi.NewResp(api.AutoRuleForm{
		IsOpen:    rule.IsOpen,
		Match:     rule.Match,
		MatchType: rule.MatchType,
		MessageId: rule.MessageId,
		Name:      rule.Name,
		ReplyType: rule.ReplyType,
		Scenes:    rule.Scenes,
		Sort:      rule.Sort,
	})
	return
}

func (c cAutoRule) Update(ctx context.Context, req *api.AutoRuleUpdateReq) (res *baseApi.NilRes, err error) {
	rule, err := service.AutoRule().First(ctx, do.CustomerChatAutoRules{
		IsSystem:   0,
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
		Id:         ghttp.RequestFromCtx(ctx).GetRouter("id").Val(),
	})
	if err != nil {
		return
	}
	updateData := service.AutoRule().Form2Do(req.AutoRuleForm)
	_, err = service.AutoRule().Update(ctx, do.CustomerChatAutoRules{Id: rule.Id}, updateData)
	if err != nil {
		return
	}
	return baseApi.NewNilResp(), nil
}

func (c cAutoRule) Store(ctx context.Context, req *api.AutoRuleStoreReq) (res *baseApi.NilRes, err error) {
	rule := service.AutoRule().Form2Do(req.AutoRuleForm)
	rule.CustomerId = service.AdminCtx().GetCustomerId(ctx)
	_, err = service.AutoRule().Save(ctx, rule)
	if err != nil {
		return
	}

	return baseApi.NewNilResp(), err
}
func (c cAutoRule) Delete(ctx context.Context, _ *api.AutoRuleDeleteReq) (resp *baseApi.NilRes, err error) {
	rule, err := service.AutoRule().First(ctx, do.CustomerChatAutoRules{
		Id:         ghttp.RequestFromCtx(ctx).GetRouter("id"),
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
	})
	if err != nil {
		return
	}
	err = service.AutoRule().Delete(ctx, rule.Id)
	if err != nil {
		return
	}
	return baseApi.NewNilResp(), nil
}
