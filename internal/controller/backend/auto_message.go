package backend

import (
	"context"
	baseApi "gf-chat/api"
	api "gf-chat/api/backend/v1"
	"gf-chat/internal/consts"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/gogf/gf/v2/frame/g"
)

var CAutoMessage = &cAutoMessage{}

type cAutoMessage struct {
}

func (c *cAutoMessage) Index(ctx context.Context, req *api.AutoMessageListReq) (res *baseApi.ListRes[*api.AutoMessage], err error) {
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
	items, err := service.AutoMessage().ToApis(ctx, p.Items)
	res = baseApi.NewListResp(items, p.Total)
	return
}

func (c *cAutoMessage) Form(ctx context.Context, _ *api.AutoMessageFormReq) (res *baseApi.NormalRes[api.AutoMessageForm], err error) {
	id := g.RequestFromCtx(ctx).GetRouter("id").Val()
	message, err := service.AutoMessage().First(ctx, do.CustomerChatAutoMessages{
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
		Id:         id,
	})
	if err != nil {
		return
	}
	form := api.AutoMessageForm{
		Name: message.Name,
		Type: message.Type,
	}
	apiMessage, err := service.AutoMessage().ToApi(ctx, message, nil)
	if err != nil {
		return
	}
	if service.ChatMessage().IsFileType(message.Type) {
		form.File = apiMessage.File
	} else {
		switch message.Type {
		case consts.MessageTypeText:
			form.Content = apiMessage.Content
		case consts.MessageTypeNavigate:
			form.Navigator = apiMessage.Navigator
		default:
		}
	}
	return baseApi.NewResp(form), nil
}

func (c *cAutoMessage) Update(ctx context.Context, req *api.AutoMessageUpdateReq) (res *baseApi.NilRes, err error) {
	id := g.RequestFromCtx(ctx).GetRouter("id").Val()
	_, err = service.AutoMessage().First(ctx, do.CustomerChatAutoMessages{
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
		Id:         id,
	})
	if err != nil {
		return
	}
	_, err = service.AutoMessage().Update(ctx, do.CustomerChatAutoMessages{Id: id},
		service.AutoMessage().Form2Do(req.AutoMessageForm))
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
	err = service.AutoMessage().Delete(ctx, message.Id)
	if err != nil {
		return
	}
	return baseApi.NewNilResp(), nil
}

func (c *cAutoMessage) Store(ctx context.Context, req *api.AutoMessageStoreReq) (res *baseApi.NilRes, err error) {
	message := service.AutoMessage().Form2Do(req.AutoMessageForm)
	message.CustomerId = service.AdminCtx().GetCustomerId(ctx)
	_, err = service.AutoMessage().Save(ctx, message)
	if err != nil {
		return
	}
	return baseApi.NewNilResp(), nil
}
