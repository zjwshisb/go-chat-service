package backend

import (
	"context"
	baseApi "gf-chat/api"
	api "gf-chat/api/v1/backend/autorule"
	"gf-chat/internal/consts"
	"gf-chat/internal/dao"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

var CAutoRule = &cAutoRule{}

type cAutoRule struct {
}

func (c cAutoRule) Delete(ctx context.Context, req *api.DeleteReq) (resp *baseApi.NilRes, err error) {
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
	_, _ = dao.CustomerChatAutoRules.Ctx(ctx).
		WherePri("id", id).Delete()
	return &baseApi.NilRes{}, nil
}

func (c cAutoRule) Message(ctx context.Context, req *api.OptionMessageReq) (res *baseApi.OptionRes, err error) {
	items, err := service.AutoMessage().All(ctx, do.CustomerChatAutoMessages{
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
	})
	if err != nil {
		return
	}
	r := baseApi.OptionRes{}
	slice.ForEach(items, func(index int, item *entity.CustomerChatAutoMessages) {
		r = append(r, model.Option{
			Value: item.Id,
			Label: item.Name,
		})
	})
	res = &r
	return
}
func (c cAutoRule) Form(ctx context.Context, req *api.FormReq) (res *api.FormRes, err error) {
	rule, err := service.AutoRule().First(ctx, do.CustomerChatAutoRules{
		IsSystem:   0,
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
		Id:         ghttp.RequestFromCtx(ctx).GetRouter("id").Val(),
	})
	if rule == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound)
	}
	scenes := make([]entity.CustomerChatAutoRuleScenes, 0)
	dao.CustomerChatAutoRuleScenes.Ctx(ctx).Where(do.CustomerChatAutoRuleScenes{
		RuleId: rule.Id,
	}).Scan(&scenes)
	sceneStr := slice.Map(scenes, func(index int, item entity.CustomerChatAutoRuleScenes) string {
		return item.Name
	})
	res = &api.FormRes{
		IsOpen:    gconv.Bool(rule.IsOpen),
		Match:     rule.Match,
		MatchType: rule.MatchType,
		MessageId: rule.MessageId,
		Name:      rule.Name,
		ReplyType: rule.ReplyType,
		Scenes:    sceneStr,
		Sort:      rule.Sort,
	}
	return
}

func (c cAutoRule) Update(ctx context.Context, req *api.UpdateReq) (res *baseApi.NilRes, err error) {
	rule, err := service.AutoRule().First(ctx, do.CustomerChatAutoRules{
		IsSystem:   0,
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
		Id:         ghttp.RequestFromCtx(ctx).GetRouter("id").Val(),
	})
	if rule == nil {
		return nil, err
	}
	err = dao.CustomerChatAutoRules.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		rule.Name = req.Name
		rule.Sort = req.Sort
		rule.IsOpen = gconv.Int(req.IsOpen)
		rule.Match = req.Match
		rule.MatchType = req.MatchType
		rule.ReplyType = req.ReplyType
		if rule.ReplyType == consts.AutoRuleReplyTypeMessage {
			rule.MessageId = req.MessageId
		}
		dao.CustomerChatAutoRules.Ctx(ctx).Save(rule)
		dao.CustomerChatAutoRuleScenes.Ctx(ctx).Where(do.CustomerChatAutoRuleScenes{RuleId: rule.Id}).Delete()
		if rule.ReplyType == consts.AutoRuleReplyTypeTransfer {
			req.Scenes = []string{consts.AutoRuleSceneNotAccepted}
		}
		for _, s := range req.Scenes {
			dao.CustomerChatAutoRuleScenes.Ctx(ctx).Data(&entity.CustomerChatAutoRuleScenes{
				Name:      s,
				RuleId:    rule.Id,
				CreatedAt: gtime.New(),
				UpdatedAt: gtime.New(),
			}).Save()
		}
		return nil
	})
	if err != nil {
		return
	}
	return &baseApi.NilRes{}, nil
}

func (c cAutoRule) Store(ctx context.Context, req *api.StoreReq) (res *baseApi.NilRes, err error) {
	err = dao.CustomerChatAutoRules.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		open := 0
		if req.IsOpen {
			open = 1
		}
		rule := &entity.CustomerChatAutoRules{
			Name:       req.Name,
			Match:      req.Match,
			MatchType:  req.MatchType,
			ReplyType:  req.ReplyType,
			IsOpen:     open,
			Sort:       req.Sort,
			IsSystem:   0,
			CustomerId: service.AdminCtx().GetCustomerId(ctx),
		}
		if rule.ReplyType == consts.AutoRuleReplyTypeMessage {
			rule.MessageId = req.MessageId
		}
		result, _ := dao.CustomerChatAutoRules.Ctx(ctx).Save(rule)
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		if rule.ReplyType == consts.AutoRuleReplyTypeTransfer {
			req.Scenes = []string{consts.AutoRuleSceneNotAccepted}
		}
		for _, s := range req.Scenes {
			dao.CustomerChatAutoRuleScenes.Ctx(ctx).Data(&entity.CustomerChatAutoRuleScenes{
				Name:      s,
				RuleId:    uint(id),
				CreatedAt: gtime.New(),
				UpdatedAt: gtime.New(),
			}).Save()
		}
		return nil
	})
	return &baseApi.NilRes{}, err
}

func (c cAutoRule) Index(ctx context.Context, req *api.ListReq) (res *api.ListRes, err error) {
	paginator, err := service.AutoRule().Paginate(ctx, &do.CustomerChatAutoRules{
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
		IsSystem:   0,
	}, model.QueryInput{
		Size: req.PageSize,
		Page: req.Current,
	}, nil, nil)
	if err != nil {
		return
	}
	items := make([]api.ListItem, len(paginator.Items))

	var messages []entity.CustomerChatAutoMessages
	messageIds := slice.Map(paginator.Items, func(index int, item model.CustomerChatAutoRule) uint {
		return item.MessageId
	})

	_ = dao.CustomerChatAutoMessages.Ctx(ctx).
		WhereIn("id", messageIds).Scan(&messages)
	for index, i := range paginator.Items {
		items[index] = api.ListItem{
			Id:        i.Id,
			Name:      i.Name,
			Match:     i.Match,
			MatchType: i.MatchType,
			ReplyType: i.ReplyType,
			IsOpen:    i.IsOpen == 1,
			Sort:      i.Sort,
			Count:     uint(i.Count),
		}
		if i.MessageId != 0 {
			m, exit := slice.Find(messages, func(index int, item entity.CustomerChatAutoMessages) bool {
				return item.Id == i.MessageId
			})
			if exit {
				items[index].Message = service.AutoMessage().EntityToListItem(*m)
			}
		}

	}

	return &api.ListRes{
		Items: items,
		Total: paginator.Total,
	}, nil
}
