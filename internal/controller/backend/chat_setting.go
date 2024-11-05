package backend

import (
	"context"
	"database/sql"
	"encoding/json"
	baseApi "gf-chat/api"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/consts"
	"gf-chat/internal/dao"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

var CChatSetting = &cChatSetting{}

type cChatSetting struct {
}

func (c *cChatSetting) Index(ctx context.Context, req *api.ChatSettingIndexReq) (res *api.ChatSettingIndexRes, err error) {
	customerId := service.AdminCtx().GetCustomerId(ctx)
	var items []entity.CustomerChatSettings
	dao.CustomerChatSettings.Ctx(ctx).Where(do.CustomerChatSettings{
		CustomerId: customerId,
	}).Scan(&items)
	r := api.ChatSettingIndexRes{}
	for _, i := range items {
		var o = make([]model.Option, 0)
		switch i.Name {
		case consts.ChatSettingOfflineSmsId:
			for _, s := range service.Sms().GetValidTemplate(customerId) {
				o = append(o, model.Option{
					Value: gconv.String(s.Id),
					Label: s.Name,
				})
			}
		case consts.ChatSettingOfflineTmplId:
			for _, s := range service.SubscribeMsg().GetEntities(customerId) {
				o = append(o, model.Option{
					Value: gconv.String(s.Id),
					Label: s.Title,
				})
			}
		default:
			_ = json.Unmarshal([]byte(i.Options), &o)
		}
		var value any
		value = i.Value
		if i.Type == "image" {
			value = service.Qiniu().Form(i.Value)
		}
		r = append(r, api.ChatSettingListItem{
			Id:          i.Id,
			Name:        i.Name,
			Value:       value,
			Options:     o,
			Title:       i.Title,
			Type:        i.Type,
			Description: i.Description,
		})
	}
	res = &r
	return res, err
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
	switch setting.Name {
	case consts.ChatSettingOfflineTmplId:
		tmpl := service.SubscribeMsg().First(do.WeappSubscribeMessages{
			Id:         req.Value,
			CustomerId: customerId,
		})
		if tmpl == nil {
			return nil, gerror.NewCode(gcode.CodeNotFound)
		}
		err = service.SubscribeMsg().CheckChatTmpl(*tmpl)
		if err != nil {
			return nil, gerror.NewCode(gcode.CodeValidationFailed, err.Error())
		}
	case consts.ChatSettingOfflineSmsId:
		sms := service.Sms().First(do.SmsTemplates{
			Id:         req.Value,
			CustomerId: customerId,
		})
		if sms == nil {
			return nil, gerror.NewCode(gcode.CodeNotFound)
		}
		err = service.Sms().CheckChatSms(sms)
		if err != nil {
			return nil, gerror.NewCode(gcode.CodeValidationFailed, err.Error())
		}
	}
	setting.Value = req.Value
	dao.CustomerChatSettings.Ctx(ctx).Save(setting)
	return &baseApi.NilRes{}, nil
}
