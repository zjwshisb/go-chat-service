package backend

import (
	"context"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/service"
)

var CDashboard = &cDashboard{}

type cDashboard struct {
}

func (c cDashboard) OnlineAdmin(ctx context.Context, req *api.DashboardOnlineAdminReq) (*api.DashboardOnlineUserRes, error) {
	res := api.DashboardOnlineUserRes{}
	user := service.Chat().GetOnlineAdmin(service.AdminCtx().GetCustomerId(ctx))
	for _, u := range user {
		res = append(res, u)
	}
	return &res, nil
}

func (c cDashboard) OnlineUser(ctx context.Context, req *api.DashboardOnlineUserReq) (*api.DashboardOnlineUserRes, error) {
	res := api.DashboardOnlineUserRes{}
	user := service.Chat().GetOnlineUser(service.AdminCtx().GetCustomerId(ctx))
	for _, u := range user {
		res = append(res, u)
	}
	return &res, nil
}

func (c cDashboard) OnlineInfo(ctx context.Context, req *api.DashboardOnlineReq) (*api.DashboardOnlineRes, error) {
	count := service.Chat().GetOnlineCount(ctx, service.AdminCtx().GetCustomerId(ctx))
	return &api.DashboardOnlineRes{
		UserCount:        count.User,
		WaitingUserCount: count.Waiting,
		AdminCount:       count.Admin,
	}, nil
}
