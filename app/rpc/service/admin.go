package service

import (
	"context"
	"ws/app/http/websocket"
	"ws/app/repositories"
	"ws/app/rpc/request"
	"ws/app/rpc/response"
)

type Admin struct {
}

func (admin *Admin) UpdateSetting(ctx context.Context, request *request.IdRequest, response *response.NilResponse) error {
	u := repositories.AdminRepo.FirstById(request.Id)
	if u != nil {
		websocket.AdminManager.UpdateSetting(u)
	}
	return nil
}

func (admin *Admin) UserTransfer(ctx context.Context, request *request.IdRequest, response *response.NilResponse) error {
	u := repositories.AdminRepo.FirstById(request.Id)
	if admin != nil {
		websocket.AdminManager.NoticeLocalUserTransfer(u)
	}
	return nil
}

func (admin *Admin) WaitingUser(ctx context.Context, request *request.GroupRequest, response *response.NilResponse) error {
	websocket.AdminManager.BroadcastLocalWaitingUser(request.GroupId)
	return nil
}

func (admin *Admin) UserOffline(ctx context.Context, request *request.IdRequest, response *response.NilResponse) error {
	websocket.AdminManager.NoticeLocalUserOffline(request.Id)
	return nil
}

func (admin *Admin) UserOnline(ctx context.Context, request *request.IdRequest, response *response.NilResponse) error {
	websocket.AdminManager.NoticeLocalUserOnline(request.Id)
	return nil
}

func (admin *Admin) OnlineAdmin(ctx context.Context, request *request.GroupRequest, response *response.NilResponse) error {
	websocket.AdminManager.BroadcastLocalOnlineAdmins(request.GroupId)
	return nil
}
