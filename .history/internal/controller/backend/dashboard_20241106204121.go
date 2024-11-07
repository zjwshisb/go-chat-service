package backend

import (
	"context"
	api "gf-chat/api/v1/backend/dashboard"
	"gf-chat/internal/service"
)

var CDashboard = &cDashboard{}

type cDashboard struct {
}

func (c cDashboard) OnlineAdmin(ctx context.Context, req *api.OnlineAdminReq) (*api.OnlineUserRes, error) {
	res := api.OnlineUserRes{}
	user := service.Chat().GetOnlineAdmin(service.AdminCtx().GetCustomerId(ctx))
	for _, u := range user {
		res = append(res, u)
	}
	return &res, nil
}

func (c cDashboard) OnlineUser(ctx context.Context, req *api.OnlineUserReq) (*api.OnlineUserRes, error) {
	res := api.OnlineUserRes{}
	user := service.Chat().GetOnlineUser(service.AdminCtx().GetCustomerId(ctx))
	for _, u := range user {
		res = append(res, u)
	}
	return &res, nil
}

func (c cDashboard) OnlineInfo(ctx context.Context, req *api.OnlineReq) (*api.OnlineRes, error) {
	count := service.Chat().GetOnlineCount(service.AdminCtx().GetCustomerId(ctx))
	return &api.OnlineRes{
		UserCount:        count.User,
		WaitingUserCount: count.Waiting,
		AdminCount:       count.Admin,
	}, nil
}
