package backend

import (
	"context"
	baseApi "gf-chat/api"
	api "gf-chat/api/v1/backend/session"
	"gf-chat/internal/consts"
	"gf-chat/internal/model"
	"gf-chat/internal/model/chat"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
)

var CSession = &cSession{}

type cSession struct {
}

func (c cSession) Index(ctx context.Context, req *api.ListReq) (*api.ListRes, error) {
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
		users := service.User().GetUsers(ctx, do.Users{
			Username:   req.Username,
			CustomerId: customerId,
		})
		uids := slice.Map(users, func(index int, item *entity.Users) uint {
			return item.Id
		})
		w["user_id"] = uids
	}
	if req.AdminName != "" {
		admins := service.Admin().GetAdmins(ctx, do.CustomerAdmins{
			Username:   req.AdminName,
			CustomerId: customerId,
		})
		adminIds := slice.Map(admins, func(index int, item *entity.CustomerAdmins) uint {
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
	session, total := service.ChatSession().Paginate(ctx, w, model.QueryInput{
		Size:      req.PageSize,
		Page:      req.Current,
		WithTotal: true,
	})
	res := make([]chat.Session, len(session), len(session))
	for index, s := range session {
		res[index] = service.ChatSession().RelationToChat(s)
	}
	r := &api.ListRes{
		Items: res,
		Total: total,
	}
	return r, nil
}

func (c cSession) Cancel(ctx context.Context, req *api.CancelReq) (*baseApi.NilRes, error) {
	id := ghttp.RequestFromCtx(ctx).GetRouter("id")
	session := service.ChatSession().First(ctx, do.CustomerChatSessions{
		Id:         id,
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
	})
	if session == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound)
	}
	err := service.ChatSession().Cancel(session)
	if err != nil {
		return nil, err
	}
	return &baseApi.NilRes{}, nil
}

func (c cSession) Close(ctx context.Context, req *api.CloseReq) (*baseApi.NilRes, error) {
	id := ghttp.RequestFromCtx(ctx).GetRouter("id")
	session := service.ChatSession().First(ctx, do.CustomerChatSessions{
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
		Id:         id,
	})
	if session == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound)
	}
	if session.BrokenAt != nil {
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "该会话已关闭")
	}
	service.ChatSession().Close(session, false, true)
	return &baseApi.NilRes{}, nil
}

func (c cSession) Detail(ctx context.Context, req *api.DetailReq) (res *api.DetailRes, err error) {
	id := ghttp.RequestFromCtx(ctx).GetRouter("id")
	session := service.ChatSession().FirstRelation(ctx, do.CustomerChatSessions{
		Id:         id,
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
	})
	if session == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound)
	}
	relations := service.ChatMessage().GetModels(0, do.CustomerChatMessages{
		SessionId: session.Id,
		Source:    []int{consts.MessageSourceAdmin, consts.MessageSourceUser},
	}, 0)
	message := make([]chat.Message, len(relations), len(relations))
	for index, i := range relations {
		message[index] = service.ChatMessage().RelationToChat(*i)
	}
	return &api.DetailRes{
		Messages: message,
		Session:  service.ChatSession().RelationToChat(session),
		Total:    len(message),
	}, nil
}
