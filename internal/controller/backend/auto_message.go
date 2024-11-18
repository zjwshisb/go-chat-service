package backend

import (
	"context"
	"encoding/json"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/consts"
	"gf-chat/internal/dao"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/frame/g"

	baseApi "gf-chat/api"

	"gf-chat/internal/model"
	"gf-chat/internal/service"
)

var CAutoMessage = &cAutoMessage{}

type cAutoMessage struct {
}

func (c *cAutoMessage) Index(ctx context.Context, req *api.AutoMessageListReq) (res *baseApi.ListRes[api.AutoMessageListItem], err error) {
	w := do.CustomerChatAutoMessages{
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
	}
	if req.Type != "" {
		w.Type = req.Type
	}
	p, err := service.AutoMessage().Paginate(ctx, &w,
		model.QueryInput{
			Page: req.Current,
			Size: req.PageSize,
		}, nil, nil)
	if err != nil {
		return
	}
	items := make([]api.AutoMessageListItem, p.Total)
	for index, i := range p.Items {
		item := api.AutoMessageListItem{
			Id:         i.Id,
			Name:       i.Name,
			Type:       i.Type,
			Content:    i.Content,
			CreatedAt:  i.CreatedAt,
			UpdatedAt:  i.UpdatedAt,
			RulesCount: 0,
		}
		if i.Type == consts.MessageTypeImage {
			//l.Content = service.Qiniu().Url(i.Content)
		}
		if i.Type == consts.MessageTypeNavigate {
			m := make(map[string]string)
			_ = json.Unmarshal([]byte(i.Content), &m)
			item.Title = m["title"]
			item.Content = m["content"]
			item.Url = m["url"]
		}
		items[index] = item
	}
	res = baseApi.NewListResp(items, p.Total)
	return
}

func (c *cAutoMessage) Form(ctx context.Context, req *api.AutoMessageFormReq) (res *baseApi.NormalRes[api.AutoMessageFormRes], err error) {
	id := g.RequestFromCtx(ctx).GetRouter("id").Val()
	message, err := service.AutoMessage().First(ctx, do.CustomerChatAutoMessages{
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
		Id:         id,
	})
	if err != nil {
		return
	}
	form := api.AutoMessageFormRes{
		Name:    message.Name,
		Type:    message.Type,
		Content: "",
		Title:   "",
		Url:     "",
	}
	switch message.Type {
	case consts.MessageTypeNavigate:
		var data map[string]string
		e := json.Unmarshal([]byte(message.Content), &data)
		if e == nil {
			//form.Content = service.Qiniu().Form(data["content"])
			form.Url = data["url"]
			form.Title = data["title"]
		}
	case consts.MessageTypeImage:
		//form.Content = service.Qiniu().Form(message.Content)
	default:
		form.Content = message.Content
	}
	return baseApi.NewResp(form), nil
}

func (c *cAutoMessage) Update(ctx context.Context, req *api.AutoMessageUpdateReq) (res *baseApi.NilRes, err error) {
	id := g.RequestFromCtx(ctx).GetRouter("id").Val()
	message, err := service.AutoMessage().First(ctx, do.CustomerChatAutoMessages{
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
		Id:         id,
	})
	if err != nil {
		return
	}
	message.Name = req.Name
	switch message.Type {
	case consts.MessageTypeNavigate:
		content := map[string]string{
			"title":   req.Title,
			"url":     req.Url,
			"content": req.Content,
		}
		contentJson, _ := json.Marshal(content)
		message.Content = string(contentJson)
	default:
		message.Content = req.Content
	}
	_, err = dao.CustomerChatAutoMessages.Ctx(ctx).Save(message)
	if err != nil {
		return
	}
	return baseApi.NewNilResp(), nil
}

func (c *cAutoMessage) Delete(ctx context.Context, req *api.AutoMessageDeleteReq) (res *baseApi.NilRes, err error) {
	id := g.RequestFromCtx(ctx).GetRouter("id").Val()
	message, err := service.AutoMessage().First(ctx, do.CustomerChatAutoMessages{
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
		Id:         id,
	})
	if err != nil {
		return
	}
	_, err = dao.CustomerChatAutoMessages.Ctx(ctx).Where("id", message.Id).Delete()
	if err != nil {
		return
	}
	return baseApi.NewNilResp(), nil
}

func (c *cAutoMessage) Store(ctx context.Context, req *api.AutoMessageStoreReq) (res *baseApi.NilRes, err error) {
	admin := service.AdminCtx().GetAdmin(ctx)
	item := entity.CustomerChatAutoMessages{
		Name:       req.Name,
		Type:       req.Type,
		CustomerId: admin.CustomerId,
	}
	switch item.Type {
	case consts.MessageTypeNavigate:
		content := map[string]string{
			"title":   req.Title,
			"url":     req.Url,
			"content": req.Content,
		}
		contentJson, _ := json.Marshal(content)
		item.Content = string(contentJson)
	default:
		item.Content = req.Content
	}
	_, err = service.AutoMessage().Insert(ctx, &model.CustomerChatAutoMessage{
		CustomerChatAutoMessages: item,
	})
	if err != nil {
		return
	}
	return baseApi.NewNilResp(), nil
}

func (c *cAutoMessage) Option(ctx context.Context, req *api.AutoMessageOptionReq) (res *baseApi.OptionRes, err error) {
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
