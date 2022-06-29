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

type Connection struct {
}

func (connection *Connection) Count(ctx context.Context, request *request.NormalRequest,
	response *response.CountResponse) error {
	if request.Types == websocket.TypeAdmin {
		response.Data = websocket.AdminManager.GetLocalOnlineTotal(request.GroupId)
	} else {
		response.Data = websocket.UserManager.GetLocalOnlineTotal(request.GroupId)
	}
	return nil
}

func (connection *Connection) AllCount(
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

func (connection *Connection) Ids(ctx context.Context, request *request.NormalRequest, response *response.IdsResponse) error {
	var ids []int64
	if request.Types == websocket.TypeUser {
		ids = websocket.UserManager.GetLocalOnlineUserIds(request.GroupId)
	} else {
		ids = websocket.AdminManager.GetLocalOnlineUserIds(request.GroupId)
	}
	response.Data = ids
	return nil
}

func (connection *Connection) Online(ctx context.Context, request *request.OnlineRequest, response *response.OnlineResponse) error {
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
	response.Data = m.IsLocalOnline(user)
	return nil
}

func (connection *Connection) RepeatConnect(ctx context.Context, request *request.RepeatConnectRequest, response *response.NilResponse) error {
	if request.Types == websocket.TypeAdmin {
		admin := repositories.AdminRepo.FirstById(request.Id)
		if admin != nil {
			websocket.AdminManager.NoticeLocalRepeatConnect(admin, request.NewUuid)
		}
	}
	return nil
}
