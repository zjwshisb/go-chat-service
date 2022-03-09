package connection

import (
	"context"
	"ws/app/repositories"
	"ws/app/rpc/request"
	"ws/app/rpc/response"
	"ws/app/websocket"
)


type Connection struct {}

func (conn *Connection) Exist(ctx context.Context, request *request.ExistRequest, response *response.ExistResponse) (err error) {
	repo := repositories.UserRepo{}
	user := repo.First([]*repositories.Where{
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

func (conn *Connection) Total(ctx context.Context, request *request.GroupRequest, response *response.TotalResponse) (err error) {
	if request.Types == "admin" {
		response.Total = websocket.AdminManager.GetLocalOnlineTotal(request.GroupId)
	} else {
		response.Total = websocket.UserManager.GetLocalOnlineTotal(request.GroupId)
	}
	return nil
}



