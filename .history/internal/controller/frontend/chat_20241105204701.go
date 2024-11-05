package frontend

import (
	"context"
	baseApi "gf-chat/api"
	"gf-chat/api/v1/frontend"
	"gf-chat/internal/consts"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
)

var CChat = &cChat{}

type cChat struct {
}

func (c cChat) Message(ctx context.Context, req *frontend.ChatMessageReq) (res *frontend.ChatMessageRes, err error) {
	uid := service.UserCtx().GetUser(ctx).Id
	messages := service.ChatMessage().GetModels(req.Id, do.CustomerChatMessages{
		UserId: uid,
	}, 20)
	r := frontend.ChatMessageRes{}
	adminToMessageId := make(map[uint][]uint)
	for _, item := range messages {
		r = append(r, service.ChatMessage().RelationToChat(*item))
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
			service.ChatMessage().ChangeToRead(ids)
			if aId > 0 {
				service.Chat().NoticeUserRead(service.UserCtx().GetCustomerId(ctx), aId, ids)
			}
		}

	}()
	return &r, nil
}

func (c cChat) Read(ctx context.Context, req *frontend.ChatReadReq) (res *baseApi.NilRes, err error) {
	user := service.UserCtx().GetUser(ctx)
	msgIds := []int64{req.MsgId}
	message := service.ChatMessage().First(do.CustomerChatMessages{
		Id:     req.MsgId,
		UserId: user.Id,
		Source: consts.MessageSourceAdmin,
	})
	if message == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound)
	}
	if message.SendAt == 0 {
		service.ChatMessage().ChangeToRead([]int64{message.Id})
		service.Chat().NoticeUserRead(user.CustomerId, message.AdminId, msgIds)
	}
	return &baseApi.NilRes{}, nil
}

func (c cChat) Rate(ctx context.Context, req *frontend.ChatRateReq) (res *baseApi.NilRes, err error) {
	id := ghttp.RequestFromCtx(ctx).GetRouter("id")
	msg := service.ChatMessage().First(do.CustomerChatMessages{
		Id:     id,
		UserId: service.UserCtx().GetUser(ctx).Id,
		Type:   consts.MessageTypeRate,
	})
	if msg == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound)
	}
	session := service.ChatSession().First(ctx, do.CustomerChatSessions{Id: msg.SessionId})
	if session != nil {
		msg.Content = gconv.String(req.Rate)
		service.ChatMessage().SaveOne(msg)
		session.Rate = req.Rate
		service.ChatSession().SaveEntity(session)
		service.Chat().NoticeRate(msg)
	}

	return &baseApi.NilRes{}, nil
}
