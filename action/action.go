package action

import (
	"encoding/json"
	"errors"
	"github.com/mitchellh/mapstructure"
	"time"
	"ws/models"
)

const (
	ReceiptAction = "receipt"
	PingAction = "ping"
	UserOnLineAction = "user-online"
	UserOffLineAction = "user-offline"
	UserWaitingCountAction = "waiting-users"
	ServerUserListAction = "server-user-list"
	MessageAction = "message"
)

type Action struct {
	Data map[string]interface{} `json:"data"`
	Time int64	`json:"time"`
	Action string `json:"action"`
	Message *models.Message `json:"-"`
}

func (action *Action) Marshal() (b []byte, err error) {
	if action.Action == PingAction {
		return []byte(""), nil
	}
	b, err = json.Marshal(action)
	return
}

func (action *Action) UnMarshal(b []byte) (err error) {
	err =  json.Unmarshal(b, action)
	return
}

func (action *Action) GetMessage() (message *models.Message,err error)  {
	if action.Action == MessageAction {
		message = &models.Message{}
		err = mapstructure.Decode(action.Data, message)
	} else {
		err = errors.New("无效的action")
	}
	return

}

func NewServerUserList(d []*models.ChatUser) *Action {
	data := make(map[string]interface{})
	data["list"] = d
	return &Action{
		Action: ServerUserListAction,
		Time: time.Now().Unix(),
		Data: data,
	}
}
func NewReceipt(action *Action) (act *Action, err error) {
	if action.Action != MessageAction {
		err = errors.New("无效的action")
		return
	}
	data := make(map[string]interface{})
	userId := action.Data["user_id"]
	data["user_id"] = userId
	data["req_id"] = action.Data["req_id"]
	act = &Action{
		Action: ReceiptAction,
		Time: time.Now().Unix(),
		Data: data,
	}
	return
}
func NewUserOnline(uid int64) *Action {
	data := make(map[string]interface{})
	data["user_id"] = uid
	return &Action{
		Action: UserOnLineAction,
		Time: time.Now().Unix(),
		Data: data,
	}
}
func NewUserOffline(uid int64) *Action {
	data := make(map[string]interface{})
	data["user_id"] = uid
	return &Action{
		Action: UserOffLineAction,
		Time: time.Now().Unix(),
		Data: data,
	}
}
func NewPing() *Action {
	return &Action{
		Action: PingAction,
		Time: time.Now().Unix(),
	}
}
func NewWaitingUsers(i interface{}) *Action {
	data := make(map[string]interface{})
	data["list"] = i
	return &Action{
		Action: UserWaitingCountAction,
		Time: time.Now().Unix(),
		Data: data,
	}
}


