package connection

import (
	"context"
	"ws/app/http/websocket"
	"ws/app/repositories"
	"ws/app/rpc/request"
	"ws/app/rpc/response"
)

type Connection struct{}

func (conn *Connection) Exist(ctx context.Context, request *request.ConnectionExistRequest,
	response *response.ConnectionExistResponse) (err error) {
	user := repositories.UserRepo.First([]*repositories.Where{
		{
			Filed: "id = ?",
			Value: request.Uid,
		},
	}, []string{})
	if user != nil {
		response.Exists = websocket.UserManager.ConnExist(user)
	}
	return err
}
func (conn *Connection) Ids(ctx context.Context,
	request *request.ConnectionGroupRequest,
	response *response.ConnectionIdsResponse) error {
	if request.Types == "admin" {
		response.Ids = websocket.AdminManager.GetLocalOnlineUserIds(request.GroupId)
	} else {
		response.Ids = websocket.UserManager.GetLocalOnlineUserIds(request.GroupId)
	}
	return nil
}

func (conn *Connection) Total(ctx context.Context, request *request.ConnectionGroupRequest, response *response.ConnectionTotalResponse) (err error) {
	if request.Types == "admin" {
		response.Total = websocket.AdminManager.GetLocalOnlineTotal(request.GroupId)
	} else {
		response.Total = websocket.UserManager.GetLocalOnlineTotal(request.GroupId)
	}
	return nil
}
func (conn *Connection) AllTotal(ctx context.Context, request *request.ConnectionGroupRequest, response *response.ConnectionTotalResponse) (err error) {
	if request.Types == "admin" {
		response.Total = websocket.AdminManager.GetAllConnCount()
	} else {
		response.Total = websocket.UserManager.GetAllConnCount()
	}
	return nil
}
