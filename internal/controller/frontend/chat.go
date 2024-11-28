package frontend

import (
	"context"
	baseApi "gf-chat/api"
	api "gf-chat/api/v1/frontend"
	"gf-chat/internal/consts"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
)

var CChat = &cChat{}

type cChat struct {
}

func (c cChat) Message(ctx context.Context, req *api.ChatMessageReq) (res *api.ChatMessageRes, err error) {
	uid := service.UserCtx().GetUser(ctx).Id
	messages, err := service.ChatMessage().GetList(ctx, req.Id, do.CustomerChatMessages{
		UserId: uid,
	}, 20)
	if err != nil {
		return
	}
	r := api.ChatMessageRes{}
	adminToMessageId := make(map[uint][]uint)
	for _, item := range messages {
		msg, err := service.ChatMessage().ToApi(ctx, item)
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
		for aId, ids := range adminToMessageId {
			_, err := service.ChatMessage().ToRead(ctx, ids)
			if err != nil {
				g.Log().Error(ctx, err)
			}
			if aId > 0 {
				service.Chat().NoticeUserRead(service.UserCtx().GetCustomerId(ctx), aId, ids)
			}
		}

	}()
	return &r, nil
}

func (c cChat) Read(ctx context.Context, req *api.ChatReadReq) (res *baseApi.NilRes, err error) {
	user := service.UserCtx().GetUser(ctx)
	msgIds := []uint{req.MsgId}
	message, err := service.ChatMessage().First(ctx, do.CustomerChatMessages{
		Id:     req.MsgId,
		UserId: user.Id,
		Source: consts.MessageSourceAdmin,
	})
	if err != nil {
		return
	}

	if message.SendAt == nil {
		_, err = service.ChatMessage().ToRead(ctx, message.Id)
		if err != nil {
			return
		}
		service.Chat().NoticeUserRead(user.CustomerId, message.AdminId, msgIds)
	}
	return &baseApi.NilRes{}, nil
}

func (c cChat) Rate(ctx context.Context, req *api.ChatRateReq) (res *baseApi.NilRes, err error) {
	id := ghttp.RequestFromCtx(ctx).GetRouter("id")
	msg, err := service.ChatMessage().First(ctx, do.CustomerChatMessages{
		Id:     id,
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
	return &baseApi.NilRes{}, nil
}
