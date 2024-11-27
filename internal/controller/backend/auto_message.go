package backend

import (
	"context"
	"encoding/json"
	baseApi "gf-chat/api"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/consts"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

var CAutoMessage = &cAutoMessage{}

type simpleNavigator struct {
	Image uint   `json:"image"`
	Url   string `json:"url"`
	Title string `json:"title"`
}

type cAutoMessage struct {
}

func (c *cAutoMessage) Index(ctx context.Context, req *api.AutoMessageListReq) (res *baseApi.ListRes[api.AutoMessageListItem], err error) {
	w := g.Map{
		"customer_id": service.AdminCtx().GetCustomerId(ctx),
	}
	if req.Type != "" {
		w["type"] = strutil.Trim(req.Type)
	}
	if req.Name != "" {
		w["name like"] = "%" + strutil.Trim(req.Name) + "%"
	}
	p, err := service.AutoMessage().Paginate(ctx, &w, req.Paginate, nil, nil)
	if err != nil {
		return
	}

	items := make([]api.AutoMessageListItem, p.Total)
	filesId := slice.Map(p.Items, func(index int, item model.CustomerChatAutoMessage) any {
		switch item.Type {
		case consts.MessageTypeFile:
			return item.Content
		case consts.MessageTypeNavigate:
			var navigator *simpleNavigator
			err := json.Unmarshal([]byte(item.Content), &navigator)
			if err != nil {
				return 0
			}
			return navigator.Image
		}
		return 0
	})
	files, err := service.File().All(ctx, do.CustomerChatFiles{
		Id: filesId,
	}, nil, nil)
	filesMap := slice.KeyBy(files, func(item *model.CustomerChatFile) uint {
		return item.Id
	})
	if err != nil {
		return
	}
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
		if i.Type == consts.MessageTypeFile {
			file := maputil.GetOrDefault(filesMap, gconv.Uint(i.Content), nil)
			if file != nil {
				item.File = service.File().ToApi(file)
			}
		}
		if i.Type == consts.MessageTypeNavigate {
			var simpleNavigator *simpleNavigator
			err = json.Unmarshal([]byte(i.Content), &simpleNavigator)
			if err == nil {
				navigator := api.AutoMessageNavigator{}
				navigator.Url = simpleNavigator.Url
				navigator.Title = simpleNavigator.Title
				image := maputil.GetOrDefault(filesMap, simpleNavigator.Image, nil)
				if image != nil {
					navigator.Image = service.File().ToApi(image)
				}
				item.Navigator = &navigator

			}

		}
		items[index] = item
	}
	res = baseApi.NewListResp(items, p.Total)
	return
}

func (c *cAutoMessage) Form(ctx context.Context, _ *api.AutoMessageFormReq) (res *baseApi.NormalRes[api.AutoMessageFormRes], err error) {
	id := g.RequestFromCtx(ctx).GetRouter("id").Val()
	message, err := service.AutoMessage().First(ctx, do.CustomerChatAutoMessages{
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
		Id:         id,
	})
	if err != nil {
		return
	}
	form := api.AutoMessageFormRes{
		AutoMessageForm: api.AutoMessageForm{
			Name: message.Name,
			Type: message.Type,
		},
	}
	switch message.Type {
	case consts.MessageTypeText:
		form.Content = message.Content
	case consts.MessageTypeNavigate:
		var navigator *simpleNavigator
		err = json.Unmarshal([]byte(message.Content), &navigator)
		if err != nil {
			return
		}
		form.Navigator = &api.AutoMessageNavigator{
			Url:   navigator.Url,
			Title: navigator.Title,
		}
		form.Navigator.Image, _ = service.File().FindAnd2Api(ctx, navigator.Image)

	case consts.MessageTypeFile:
		form.File, _ = service.File().FindAnd2Api(ctx, message.Content)

	default:
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
	service.AutoMessage().Fill(message, req.AutoMessageForm)
	_, err = service.AutoMessage().Save(ctx, message)
	if err != nil {
		return
	}
	return baseApi.NewNilResp(), nil
}

func (c *cAutoMessage) Delete(ctx context.Context, _ *api.AutoMessageDeleteReq) (res *baseApi.NilRes, err error) {
	id := g.RequestFromCtx(ctx).GetRouter("id").Val()
	message, err := service.AutoMessage().First(ctx, do.CustomerChatAutoMessages{
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
		Id:         id,
	})
	if err != nil {
		return
	}
	err = service.ChatMessage().Delete(ctx, message.Id)
	if err != nil {
		return
	}
	return baseApi.NewNilResp(), nil
}

func (c *cAutoMessage) Store(ctx context.Context, req *api.AutoMessageStoreReq) (res *baseApi.NilRes, err error) {
	admin := service.AdminCtx().GetUser(ctx)
	message := &model.CustomerChatAutoMessage{
		CustomerChatAutoMessages: entity.CustomerChatAutoMessages{
			Name:       req.Name,
			Type:       req.Type,
			CustomerId: admin.CustomerId,
		},
	}
	service.AutoMessage().Fill(message, req.AutoMessageForm)
	_, err = service.AutoMessage().Save(ctx, message)
	if err != nil {
		return
	}
	return baseApi.NewNilResp(), nil
}
