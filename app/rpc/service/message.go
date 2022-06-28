package service

import (
	"context"
	"ws/app/http/websocket"
	"ws/app/models"
	"ws/app/repositories"
	"ws/app/rpc/request"
	"ws/app/rpc/response"
)

type Message struct {
}

func (message *Message) Send(ctx context.Context, request *request.SendMessageRequest, response *response.SendMessageResponse) error {
	msg := repositories.MessageRepo.FirstById(request.Id)
	var m websocket.MessageHandle
	if msg != nil {
		switch msg.Source {
		case models.SourceUser:
			m = websocket.AdminManager
		case models.SourceAdmin:
			m = websocket.UserManager
		}
		m.DeliveryMessage(msg, false)
	}
	return nil
}
