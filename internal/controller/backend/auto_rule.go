package backend

import (
	"context"
	baseApi "gf-chat/api"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/consts"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
)

var CAutoRule = &cAutoRule{}

type cAutoRule struct {
}

func (c cAutoRule) Delete(ctx context.Context, _ *api.AutoRuleDeleteReq) (resp *baseApi.NilRes, err error) {
	id := ghttp.RequestFromCtx(ctx).GetRouter("id")
	rule, err := service.AutoRule().First(ctx, do.CustomerChatAutoRules{
		Id:         id,
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
	})
	if err != nil {
		return
	}
	if rule == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound)
	}
	err = service.AutoRule().Delete(ctx, id)
	if err != nil {
		return
	}
	return baseApi.NewNilResp(), nil
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
	updateData := do.CustomerChatAutoRules{
		Name:      req.Name,
		Sort:      req.Sort,
		IsOpen:    req.IsOpen,
		Match:     req.Match,
		MatchType: req.MatchType,
	}
	if rule.ReplyType == consts.AutoRuleReplyTypeMessage {
		updateData.MessageId = req.MessageId
		updateData.Scenes = req.Scenes
	}
	if rule.ReplyType == consts.AutoRuleReplyTypeTransfer {
		updateData.Scenes = []string{consts.AutoRuleSceneNotAccepted}
	}
	_, err = service.AutoRule().Update(ctx, do.CustomerChatAutoRules{Id: rule.Id}, updateData)
	if err != nil {
		return
	}
	return baseApi.NewNilResp(), nil
}

func (c cAutoRule) Store(ctx context.Context, req *api.AutoRuleStoreReq) (res *baseApi.NilRes, err error) {
	rule := &model.CustomerChatAutoRule{
		CustomerChatAutoRules: entity.CustomerChatAutoRules{
			Name:       req.Name,
			Match:      req.Match,
			MatchType:  req.MatchType,
			ReplyType:  req.ReplyType,
			Sort:       req.Sort,
			IsSystem:   0,
			CustomerId: service.AdminCtx().GetCustomerId(ctx),
		},
		IsOpen: req.IsOpen,
	}
	if rule.ReplyType == consts.AutoRuleReplyTypeMessage {
		rule.MessageId = req.MessageId
		rule.Scenes = req.Scenes
	}
	if rule.ReplyType == consts.AutoRuleReplyTypeTransfer {
		rule.Scenes = []string{consts.AutoRuleSceneNotAccepted}
	}
	_, err = service.AutoRule().Save(ctx, rule)
	if err != nil {
		return
	}

	return baseApi.NewNilResp(), err
}

func (c cAutoRule) Index(ctx context.Context, req *api.AutoRuleListReq) (res *baseApi.ListRes[api.AutoRule], err error) {
	w := g.Map{
		"customer_id": service.AdminCtx().GetCustomerId(ctx),
		"is_system":   0,
	}
	if req.Name != "" {
		w["name like ?"] = "%" + req.Name + "%"
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
