package backend

import (
	"context"
	"encoding/json"
	baseApi "gf-chat/api"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/consts"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/frame/g"
)

var CChatSetting = &cChatSetting{}

type cChatSetting struct {
}

func (c *cChatSetting) Index(ctx context.Context, req *api.ChatSettingListReq) (res *baseApi.ListRes[api.ChatSetting], err error) {
	w := g.Map{
		"customer_id": service.AdminCtx().GetCustomerId(ctx),
	}
	if req.Title != "" {
		w["title like ?"] = "%" + req.Title + "%"
	}
	items, err := service.ChatSetting().All(ctx, w, nil, nil)
	if err != nil {
		return
	}
	fileIds := slice.Map(items, func(index int, item *model.CustomerChatSetting) uint {
		if item.Type == consts.ChatSettingTypeImage {
			return gconv.Uint(item.Value)
		}
		return 0
	})
	files, err := service.File().All(ctx, do.CustomerChatFiles{
		Id: fileIds,
	}, nil, nil)
	if err != nil {
		return
	}
	filesMap := slice.KeyBy(files, func(item *model.CustomerChatFile) uint {
		return item.Id
	})
	settings := slice.Map(items, func(index int, i *model.CustomerChatSetting) api.ChatSetting {
		var value any
		value = i.Value
		if i.Type == consts.ChatSettingTypeImage {
			file := maputil.GetOrDefault(filesMap, gconv.Uint(i.Value), nil)
			if file != nil {
				value = service.File().ToApi(file)
			}

		}
		return api.ChatSetting{
			Id:          i.Id,
			Name:        i.Name,
			Value:       value,
			Options:     i.Options,
			Title:       i.Title,
			Type:        i.Type,
			Description: i.Description,
			CreatedAt:   i.CreatedAt,
			UpdatedAt:   i.UpdatedAt,
		}
	})
	return baseApi.NewListResp(settings, len(settings)), nil
}

func (c *cChatSetting) Update(ctx context.Context, req *api.ChatSettingUpdateReq) (res *baseApi.NilRes, err error) {
	setting, err := service.ChatSetting().First(ctx, do.CustomerChatSettings{
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
		Id:         g.RequestFromCtx(ctx).GetRouter("id").Val(),
	})
	if err != nil {
		return
	}
	updateData := do.CustomerChatSettings{}
	switch setting.Type {
	case consts.ChatSettingTypeText:
		updateData.Value = req.Value
	case consts.ChatSettingTypeSelect:
		for _, v := range setting.Options {
			if gconv.String(v.Value) == gconv.String(req.Value) {
				updateData.Value = req.Value
				break
			}
		}
		if updateData.Value == nil {
			err = gerror.New("无效的选项")
		}
	case consts.ChatSettingTypeImage:
		var apiFile *api.File
		err = json.Unmarshal(gconv.Bytes(req.Value), &apiFile)
		if err == nil {
			var file *model.CustomerChatFile
			file, err = service.File().First(ctx, do.CustomerChatFiles{
				Type:       consts.FileTypeImage,
				CustomerId: service.AdminCtx().GetCustomerId(ctx),
				Id:         apiFile.Id,
			})
			g.Dump(file)
			if err == nil {
				updateData.Value = file.Id
			} else {
				err = gerror.New("无效的图片")
			}
		}
	}
	if err != nil {
		return
	}
	_, err = service.ChatSetting().Update(ctx, do.CustomerChatSettings{Id: setting.Id}, updateData)
	if err != nil {
		return
	}
	return baseApi.NewNilResp(), nil
}
