package backend

import (
	"context"
	baseApi "gf-chat/api"
	api "gf-chat/api/backend/v1"
	"gf-chat/internal/consts"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"
	"github.com/duke-git/lancet/v2/slice"
)

var COption = &cOption{}

type cOption struct {
}

func (c *cOption) AutoMessage(ctx context.Context, _ *api.OptionAutoMessageReq) (res *baseApi.OptionRes, err error) {
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

func (c *cOption) MessageType(_ context.Context, _ *api.OptionMessageTypeReq) (res *baseApi.OptionRes, err error) {
	options := []baseApi.Option{
		{
			Label: "文本",
			Value: consts.MessageTypeText,
		},
		{
			Label: "图片",
			Value: consts.MessageTypeImage,
		},
		{
			Label: "音频",
			Value: consts.MessageTypeAudio,
		},
		{
			Label: "视频",
			Value: consts.MessageTypeVideo,
		},
		{
			Label: "导航卡片",
			Value: consts.MessageTypeNavigate,
		},
	}
	return baseApi.NewOptionResp(options), nil
}

func (c *cOption) AutoRuleMatchType(_ context.Context, _ *api.OptionAutoRuleMatchTypeReq) (res *baseApi.OptionRes, err error) {
	options := []baseApi.Option{
		{
			Label: "全匹配",
			Value: consts.AutoRuleMatchTypeAll,
		},
		{
			Label: "半匹配",
			Value: consts.AutoRuleMatchTypePart,
		},
	}
	return baseApi.NewOptionResp(options), nil
}

func (c *cOption) AutoRuleReplyType(_ context.Context, _ *api.OptionAutoRuleReplyTypeReq) (res *baseApi.OptionRes, err error) {
	options := []baseApi.Option{
		{
			Label: "回复消息",
			Value: consts.AutoRuleReplyTypeMessage,
		},
		{
			Label: "转接人工",
			Value: consts.AutoRuleReplyTypeTransfer,
		},
	}
	return baseApi.NewOptionResp(options), nil
}

func (c *cOption) AutoRuleScene(_ context.Context, _ *api.OptionAutoRuleSceneReq) (res *baseApi.OptionRes, err error) {
	options := []baseApi.Option{
		{
			Label: "人工未接入",
			Value: consts.AutoRuleSceneNotAccepted,
		},
		{
			Label: "已接入但客服离线",
			Value: consts.AutoRuleSceneAdminOffline,
		},
		{
			Label: "已接入客服在线",
			Value: consts.AutoRuleSceneAdminOnline,
		},
	}
	return baseApi.NewOptionResp(options), nil
}

func (c *cOption) FileType(_ context.Context, _ *api.OptionFileTypeReq) (res *baseApi.OptionRes, err error) {
	return baseApi.NewOptionResp(consts.UserAllowMessageType), nil
}

func (c *cOption) SessionStatus(_ context.Context, _ *api.OptionSessionStatusReq) (res *baseApi.OptionRes, err error) {
	options := []baseApi.Option{
		{
			Label: "已取消",
			Value: consts.ChatSessionStatusCancel,
		},
		{
			Label: "已关闭",
			Value: consts.ChatSessionStatusClose,
		},
		{
			Label: "待接入",
			Value: consts.ChatSessionStatusWait,
		},
		{
			Label: "已接入",
			Value: consts.ChatSessionStatusAccept,
		},
	}
	return baseApi.NewOptionResp(options), nil
}
