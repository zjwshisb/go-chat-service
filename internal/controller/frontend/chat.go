package frontend

import (
	"context"
	baseApi "gf-chat/api"
	"gf-chat/api/v1/backend"
	api "gf-chat/api/v1/frontend"
	"gf-chat/internal/consts"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
)

var CChat = &cChat{}

type cChat struct {
}

func (c cChat) Message(ctx context.Context, req *api.ChatMessageReq) (res *baseApi.NormalRes[[]*backend.ChatMessage], err error) {
	uid := service.UserCtx().GetUser(ctx).Id
	w := g.Map{
		"user_id": uid,
	}
	if req.Id > 0 {
		w["id <"] = req.Id
	}
	messages, err := service.ChatMessage().All(ctx, w, g.Slice{
		model.CustomerChatMessage{}.User,
		model.CustomerChatMessage{}.Admin,
	}, "id desc", req.PageSize)
	if err != nil {
		return
	}
	r := make([]*backend.ChatMessage, 0)
	adminToMessageId := make(map[uint][]uint)
	for _, item := range messages {
		msg, err := service.ChatMessage().ToApi(ctx, item, nil)
		if err != nil {
			return nil, err
		}
		r = append(r, msg)
		if item.ReadAt == nil && item.Source == consts.MessageSourceAdmin {
			ids, exist := adminToMessageId[item.AdminId]
			if exist {
				adminToMessageId[item.AdminId] = append(ids, item.Id)
			} else {
				adminToMessageId[item.AdminId] = []uint{item.Id}
			}
		}
	}
	go func() {
		for adminId, ids := range adminToMessageId {
			_, err := service.ChatMessage().ToRead(ctx, ids)
			if err != nil {
				g.Log().Error(ctx, err)
			}
			if adminId > 0 {
				service.Chat().NoticeUserRead(service.UserCtx().GetCustomerId(ctx), adminId, ids)
			}
		}

	}()
	return baseApi.NewResp(r), nil
}

func (c cChat) Read(ctx context.Context, req *api.ChatReadReq) (res *baseApi.NilRes, err error) {
	user := service.UserCtx().GetUser(ctx)
	message, err := service.ChatMessage().First(ctx, do.CustomerChatMessages{
		Id:     req.MsgId,
		UserId: user.Id,
		Source: []int{consts.MessageSourceAdmin, consts.MessageSourceSystem},
	})
	if err != nil {
		return
	}

	if message.ReadAt == nil {
		_, err = service.ChatMessage().ToRead(ctx, message.Id)
		msgIds := []uint{req.MsgId}
		if err != nil {
			return
		}
		service.Chat().NoticeUserRead(user.CustomerId, message.AdminId, msgIds)
	}
	return baseApi.NewNilResp(), nil
}

func (c cChat) Rate(ctx context.Context, req *api.ChatRateReq) (res *baseApi.NilRes, err error) {
	msg, err := service.ChatMessage().First(ctx, do.CustomerChatMessages{
		Id:     ghttp.RequestFromCtx(ctx).GetRouter("id"),
		UserId: service.UserCtx().GetUser(ctx).Id,
		Type:   consts.MessageTypeRate,
	})
	if err != nil {
		return
	}
	session, err := service.ChatSession().Find(ctx, msg.SessionId)
	if err != nil {
		return
	}
	msg.Content = gconv.String(req.Rate)
	_, err = service.ChatMessage().Save(ctx, msg)
	if err != nil {
		return
	}
	session.Rate = req.Rate
	_, err = service.ChatSession().Save(ctx, session)
	if err != nil {
		return
	}
	service.Chat().NoticeRate(msg)
	return baseApi.NewNilResp(), nil
}

func (c cChat) ReqId(_ context.Context, _ *api.ChatReqIdReq) (res *baseApi.NormalRes[api.ChatReqId], err error) {
	return baseApi.NewResp(api.ChatReqId{ReqId: service.ChatMessage().GenReqId()}), nil
}
