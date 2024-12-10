package backend

import (
	"context"
	"database/sql"
	"errors"
	baseApi "gf-chat/api"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/consts"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"
	"github.com/duke-git/lancet/v2/strutil"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

var CSession = &cSession{}

type cSession struct {
}

func (c cSession) Index(ctx context.Context, req *api.SessionListReq) (resp *baseApi.ListRes[api.ChatSession], err error) {
	w := g.Map{}
	customerId := service.AdminCtx().GetCustomerId(ctx)
	w["customer_id"] = customerId
	if startTime, exist := req.QueriedAt["0"]; exist {
		w["queried_at >="] = strutil.Trim(startTime)
	}
	if endTime, exist := req.QueriedAt["1"]; exist {
		w["queried_at <="] = strutil.Trim(endTime)
	}
	username := strutil.Trim(req.Username)
	if username != "" {
		uW := make(map[string]any)
		uW["username like"] = username + "%"
		uW["customer_id"] = customerId
		users, err := service.User().All(ctx, uW, nil, nil)
		if err != nil {
			return nil, err
		}
		uids := slice.Map(users, func(index int, item *model.User) uint {
			return item.Id
		})
		w["user_id"] = uids
	}
	adminName := strutil.Trim(req.AdminName)
	if adminName != "" {
		admins, err := service.Admin().All(ctx, g.Map{
			"username like": "%" + adminName + "%",
			"customer_id":   customerId,
		}, nil, nil)
		if err != nil {
			return nil, err
		}
		adminIds := slice.Map(admins, func(index int, item *model.CustomerAdmin) uint {
			return item.Id
		})
		w["admin_id"] = adminIds
	}
	if req.Status != "" {
		switch req.Status {
		case consts.ChatSessionStatusCancel:
			w["canceled_at is null"] = nil
		case consts.ChatSessionStatusWait:
			w["accepted_at is null"] = nil
			w["canceled_at is null"] = nil
		case consts.ChatSessionStatusAccept:
			w["accepted_at is not null"] = nil
			w["broken_at is null"] = nil
		case consts.ChatSessionStatusClose:
			w["broken_at is not null"] = nil
		}
	}
	paginator, err := service.ChatSession().Paginate(ctx, w, req.Paginate, g.Array{
		model.User{},
		model.CustomerAdmin{},
	}, "id desc")
	if err != nil {
		return
	}
	res := slice.Map(paginator.Items, func(index int, item *model.CustomerChatSession) api.ChatSession {
		return service.ChatSession().ToApi(item)
	})
	return baseApi.NewListResp(res, paginator.Total), nil
}

func (c cSession) Cancel(ctx context.Context, _ *api.SessionCancelReq) (resp *baseApi.NilRes, err error) {
	session, err := service.ChatSession().First(ctx, do.CustomerChatSessions{
		Id:         ghttp.RequestFromCtx(ctx).GetRouter("id"),
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
	})
	if err != nil {
		return
	}
	err = service.ChatSession().Cancel(ctx, session)
	if err != nil {
		return
	}
	return baseApi.NewNilResp(), nil
}

func (c cSession) Close(ctx context.Context, _ *api.SessionCloseReq) (resp *baseApi.NilRes, err error) {
	session, err := service.ChatSession().First(ctx, do.CustomerChatSessions{
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
		Id:         ghttp.RequestFromCtx(ctx).GetRouter("id"),
	})
	if err != nil {
		return
	}
	err = service.ChatSession().Close(ctx, session, false, true)
	if err != nil {
		return
	}
	return baseApi.NewNilResp(), nil
}

func (c cSession) Show(ctx context.Context, _ *api.SessionDetailReq) (res *baseApi.NormalRes[api.SessionDetail], err error) {
	session, err := service.ChatSession().First(ctx, do.CustomerChatSessions{
		Id:         ghttp.RequestFromCtx(ctx).GetRouter("id"),
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
	})

	if err != nil {
		return
	}
	if session.AdminId > 0 {
		session.Admin, err = service.Admin().Find(ctx, session.AdminId)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return
		}
	}
	session.User, err = service.User().Find(ctx, session.UserId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return
	}
	relations, err := service.ChatMessage().All(ctx, do.CustomerChatMessages{
		SessionId: session.Id,
		Source:    []int{consts.MessageSourceAdmin, consts.MessageSourceUser},
	}, nil, "id")
	if err != nil {
		return
	}
	message := make([]*api.ChatMessage, len(relations))
	for index, i := range relations {
		msg, err := service.ChatMessage().ToApi(ctx, i)
		if err != nil {
			return nil, err
		}
		message[index] = msg
	}
	return baseApi.NewResp(api.SessionDetail{
		Messages: message,
		Session:  service.ChatSession().ToApi(session),
	}), nil

}
