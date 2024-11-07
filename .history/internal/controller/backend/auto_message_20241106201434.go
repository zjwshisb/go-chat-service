package backend

import (
	"context"
	"encoding/json"
	"gf-chat/internal/consts"
	"gf-chat/internal/dao"
	"gf-chat/internal/model/do"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	baseApi "gf-chat/api"

	"gf-chat/api/v1/backend/automessage"
	"gf-chat/internal/model"
	"gf-chat/internal/service"
)

var CAutoMessage = &cAutoMessage{}

type cAutoMessage struct {
}

func (c *cAutoMessage) Index(ctx context.Context, req *automessage.ListReq) (res *automessage.ListRes, err error) {
	w := do.CustomerChatAutoMessages{
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
	}
	if req.Type != "" {
		w.Type = req.Type
	}
	entities, total := service.AutoMessage().Paginate(ctx, &w,
		model.QueryInput{
			Page: req.Current,
			Size: req.PageSize,
		})
	items := make([]model.AutoMessageListItem, len(entities))
	for index, i := range entities {
		items[index] = service.AutoMessage().EntityToListItem(*i)
	}
	res = &automessage.ListRes{
		Items: items,
		Total: total,
	}
	return
}

func (c *cAutoMessage) Form(ctx context.Context, req *automessage.FormReq) (res *automessage.FormRes, err error) {
	id := g.RequestFromCtx(ctx).GetRouter("id").Val()
	message := service.AutoMessage().First(ctx, do.CustomerChatAutoMessages{
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
		Id:         id,
	})
	if message == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound)
	}
	if err != nil {
		return nil, err
	}
	res = &automessage.FormRes{
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
			res.Content = service.Qiniu().Form(data["content"])
			res.Url = data["url"]
			res.Title = data["title"]
		}
	case consts.MessageTypeImage:
		res.Content = service.Qiniu().Form(message.Content)
	default:
		res.Content = message.Content
	}
	return
}

func (c *cAutoMessage) Update(ctx context.Context, req *automessage.UpdateReq) (res *baseApi.NilRes, err error) {
	id := g.RequestFromCtx(ctx).GetRouter("id").Val()
	message := service.AutoMessage().First(ctx, do.CustomerChatAutoMessages{
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
		Id:         id,
	})
	if message == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound)
	}
	service.AutoMessage().Update(ctx, message, req)
	return &baseApi.NilRes{}, nil
}

func (c *cAutoMessage) Delete(ctx context.Context, req *automessage.DeleteReq) (res *baseApi.NilRes, err error) {
	id := g.RequestFromCtx(ctx).GetRouter("id").Val()
	message := service.AutoMessage().First(ctx, do.CustomerChatAutoMessages{
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
		Id:         id,
	})
	if message == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound)
	}
	dao.CustomerChatAutoMessages.Ctx(ctx).Where("id", message.Id).Delete()
	return &baseApi.NilRes{}, nil
}

func (c *cAutoMessage) Store(ctx context.Context, req *automessage.StoreReq) (res *baseApi.NilRes, err error) {
	service.AutoMessage().Save(ctx, req)
	return &baseApi.NilRes{}, nil
}
