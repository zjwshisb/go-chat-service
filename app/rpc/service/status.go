package service

import (
	"context"
	"errors"
	"ws/app/contract"
	"ws/app/http/websocket"
	"ws/app/repositories"
	"ws/app/rpc/request"
	"ws/app/rpc/response"
)

type Status struct {
}

func (status *Status) count(ctx context.Context, request *request.NormalRequest,
	response *response.CountResponse) error {
	if request.Types == websocket.TypeAdmin {
		response.Data = websocket.AdminManager.GetLocalOnlineTotal(request.GroupId)
	} else {
		response.Data = websocket.UserManager.GetLocalOnlineTotal(request.GroupId)
	}
	return nil
}

func (status *Status) allCount(
	ctx context.Context,
	request *request.NormalRequest,
	response *response.CountResponse) error {
	if request.Types == websocket.TypeAdmin {
		response.Data = websocket.AdminManager.GetAllConnCount()
	} else {
		response.Data = websocket.UserManager.GetAllConnCount()
	}
	return nil
}

func (status *Status) Ids(ctx context.Context, request *request.NormalRequest, response *response.IdsResponse) error {
	var ids []int64
	if request.Types == websocket.TypeUser {
		ids = websocket.UserManager.GetLocalOnlineUserIds(request.GroupId)
	} else {
		ids = websocket.AdminManager.GetLocalOnlineUserIds(request.GroupId)
	}
	response.Data = ids
	return nil
}

func (status *Status) Online(ctx context.Context, request *request.OnlineRequest, response *response.OnlineResponse) error {
	var m websocket.ConnContainer
	var user contract.User
	if request.Types == websocket.TypeUser {
		m = websocket.UserManager
		user = repositories.UserRepo.FirstById(request.Id)
	} else {
		m = websocket.AdminManager
		user = repositories.AdminRepo.FirstById(request.Id)
	}
	if user == nil {
		return errors.New("user not exit")
	}
	response.Data = m.ConnExist(user)
	return nil
}

func (status *Status) RepeatConnect(ctx context.Context, request *request.RepeatConnectRequest, response *response.NilResponse) error {
	if request.Types == websocket.TypeAdmin {
		admin := repositories.AdminRepo.FirstById(request.Id)
		websocket.AdminManager.NoticeOtherLogin(admin)
	}
	return nil
}
