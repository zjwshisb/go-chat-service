package json

import (
	"ws/internal/models"
)

type Message struct {
	Id uint64 `json:"id"`
	UserId int64 `json:"user_id"`
	ServiceId int64 `json:"service_id"`
	Type string `json:"type"`
	Content string `json:"content"`
	ReceivedAT int64 `json:"received_at"`
	IsServer bool `json:"is_server"`
	ReqId int64 `json:"req_id"`
	IsSuccess bool `json:"is_success"`
	IsRead bool `json:"is_read"`
	Avatar string `json:"avatar"`
}

func NewMessage(model *models.Message) *Message {
	var avatar string
	if model.IsServer {
		avatar = model.BackendUser.GetAvatarUrl()
	} else {
		avatar = model.User.GetAvatarUrl()
	}
	return &Message{
		Id: model.Id,
		UserId: model.UserId,
		Type: model.Type,
		Content: model.Content,
		IsServer: model.IsServer,
		ReqId: model.ReqId,
		IsSuccess: true,
		ReceivedAT: model.ReceivedAT,
		IsRead: model.IsRead,
		Avatar: avatar,
	}
}
