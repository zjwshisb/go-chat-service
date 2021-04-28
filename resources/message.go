package resources

import (
	"fmt"
	"ws/db"
	"ws/models"
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
func NewMessage(model models.Message) *Message {
	var avatar string
	if model.IsServer {
		var serverUser models.ServerUser
		if model.ServerUser.ID == 0 {
			_ = db.Db.Model(&model).Association("ServerUser").Find(&serverUser)
		} else {
			serverUser = model.ServerUser
		}
		avatar = serverUser.GetAvatarUrl()
	} else {
		//avatar = model.User.g
	}
	fmt.Println(avatar)
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
