package backend

import (
	"context"
	baseApi "gf-chat/api"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/consts"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
)

var CSession = &cSession{}

type cSession struct {
}

func (c cSession) Index(ctx context.Context, req *api.SessionListReq) (resp *baseApi.ListRes[api.ChatSession], err error) {
	w := make(map[string]any)
	customerId := service.AdminCtx().GetCustomerId(ctx)
	w["customer_id"] = customerId
	if startTime, exist := req.QueriedAt["0"]; exist {
		w["queried_at>="] = gtime.New(startTime).Unix()
	}
	if endTime, exist := req.QueriedAt["1"]; exist {
		w["queried_at<="] = gtime.New(endTime).Unix()
	}
	if req.Username != "" {
		uW := make(map[string]any)
		uW["phone"] = req.Username
		uW["customer_id"] = customerId
		users, err := service.User().All(ctx, do.Users{
			Username:   req.Username,
			CustomerId: customerId,
		}, nil, nil)
		if err != nil {
			return nil, err
		}
		uids := slice.Map(users, func(index int, item *model.User) uint {
			return item.Id
		})
		w["user_id"] = uids
	}
	if req.AdminName != "" {
		admins, err := service.Admin().All(ctx, do.CustomerAdmins{
			Username:   req.AdminName,
			CustomerId: customerId,
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
			w["canceled_at>"] = 0
		case consts.ChatSessionStatusWait:
			w["accepted_at"] = 0
			w["canceled_at"] = 0
		case consts.ChatSessionStatusAccept:
			w["accepted_at>"] = 0
			w["broke_at"] = 0
		case consts.ChatSessionStatusClose:
			w["broke_at>"] = 0
		}
	}
	paginator, err := service.ChatSession().Paginate(ctx, w, req.Paginate, g.Array{
		model.CustomerChatSession{}.User,
		model.CustomerChatSession{}.Admin,
	}, nil)
	if err != nil {
		return
	}
	res := make([]api.ChatSession, len(paginator.Items))
	for index, s := range paginator.Items {
		res[index] = service.ChatSession().RelationToChat(s)
	}
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
	if session.BrokenAt != nil {
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "该会话已关闭")
	}
	err = service.ChatSession().Close(ctx, session, false, true)
	if err != nil {
		return
	}
	return baseApi.NewNilResp(), nil
}

func (c cSession) Detail(ctx context.Context, _ *api.SessionDetailReq) (res *api.SessionDetailRes, err error) {
	session, err := service.ChatSession().First(ctx, do.CustomerChatSessions{
		Id:         ghttp.RequestFromCtx(ctx).GetRouter("id"),
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
	})
	if err != nil {
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
	return &api.SessionDetailRes{
		Messages: message,
		Session:  service.ChatSession().RelationToChat(session),
		Total:    len(message),
	}, nil
}
