package service

import (
	"context"
	"ws/app/http/websocket"
	"ws/app/rpc/request"
	"ws/app/rpc/response"
)

type User struct {
}

func (user *User) QueueLocation(ctx context.Context, request *request.GroupRequest, response *response.NilResponse) error {
	websocket.UserManager.BroadcastLocalQueueLocation(request.GroupId)
	return nil
}
