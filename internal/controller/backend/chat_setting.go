package backend

import (
	"context"
	"database/sql"
	"encoding/json"
	baseApi "gf-chat/api"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/dao"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/frame/g"
)

var CChatSetting = &cChatSetting{}

type cChatSetting struct {
}

func (c *cChatSetting) Index(ctx context.Context, req *api.ChatSettingListReq) (res *baseApi.ListRes[api.ChatSettingListItem], err error) {
	customerId := service.AdminCtx().GetCustomerId(ctx)
	var items []entity.CustomerChatSettings
	dao.CustomerChatSettings.Ctx(ctx).Where(do.CustomerChatSettings{
		CustomerId: customerId,
	}).Scan(&items)
	settings := slice.Map(items, func(index int, i entity.CustomerChatSettings) api.ChatSettingListItem {
		var o = make([]baseApi.Option, 0)
		json.Unmarshal([]byte(i.Options), &o)
		var value any
		value = i.Value
		if i.Type == "image" {
			//value = service.Qiniu().Form(i.Value)
		}
		return api.ChatSettingListItem{
			Id:          i.Id,
			Name:        i.Name,
			Value:       value,
			Options:     o,
			Title:       i.Title,
			Type:        i.Type,
			Description: i.Description,
		}
	})
	return baseApi.NewListResp(settings, 0), nil
}

func (c *cChatSetting) Update(ctx context.Context, req *api.ChatSettingUpdateReq) (res *baseApi.NilRes, err error) {
	id := g.RequestFromCtx(ctx).GetRouter("id").Val()
	setting := &entity.CustomerChatSettings{}
	customerId := service.AdminCtx().GetCustomerId(ctx)
	err = dao.CustomerChatSettings.Ctx(ctx).Where(do.CustomerChatSettings{
		CustomerId: customerId,
		Id:         id,
	}).Scan(setting)
	if err == sql.ErrNoRows {
		return nil, err
	}
	setting.Value = req.Value
	dao.CustomerChatSettings.Ctx(ctx).Save(setting)
	return baseApi.NewNilResp(), nil
}
