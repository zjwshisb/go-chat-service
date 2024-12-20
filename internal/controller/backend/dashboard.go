package backend

import (
	"context"
	v1 "gf-chat/api/v1"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/os/gtime"
)

var CDashboard = &cDashboard{}

type cDashboard struct {
}

func (c cDashboard) AdminInfo(ctx context.Context, _ *api.DashboardAdminInfoReq) (res *v1.NormalRes[api.DashboardAdminInfo], err error) {
	user := service.Chat().GetOnlineAdmins(service.AdminCtx().GetCustomerId(ctx))
	total, err := service.Admin().Count(ctx, do.CustomerAdmins{
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
	})
	if err != nil {
		return
	}
	res = v1.NewResp(api.DashboardAdminInfo{
		Admins: user,
		Total:  total,
	})
	return
}

func (c cDashboard) WaitingUser(ctx context.Context, _ *api.DashboardWaitingUserReq) (res *v1.NormalRes[[]api.ChatSimpleUser], err error) {
	user, err := service.Chat().GetWaitingUsers(ctx, service.AdminCtx().GetCustomerId(ctx))
	if err != nil {
		return
	}
	res = v1.NewResp(user)
	return
}

func (c cDashboard) OnlineUserInfo(ctx context.Context, _ *api.DashboardOnlineUserInfoReq) (res *v1.NormalRes[api.DashboardOnlineUserInfo], err error) {
	users := service.Chat().GetOnlineUsers(service.AdminCtx().GetCustomerId(ctx))
	count, err := service.User().GetActiveCount(ctx, gtime.Now())
	res = v1.NewResp(api.DashboardOnlineUserInfo{
		Users:       users,
		ActiveCount: count,
	})
	return
}
