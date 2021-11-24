package websocket

import (
	"encoding/json"
	"errors"
	"github.com/mitchellh/mapstructure"
	"time"
	"ws/app/models"
	"ws/app/resource"
)

const (
	ReceiptAction = "receipt"
	PingAction = "ping"
	UserOnLineAction = "user-online"
	UserOffLineAction = "user-offline"
	WaitingUserAction = "waiting-users"
	WaitingUserCount = "waiting-user-count"
 	AdminsAction = "admins"
	SendMessageAction = "send-message"
	ReceiveMessageAction = "receive-message"
	OtherLogin = "other-login"
	MoreThanOne = "more-than-one"
	UserTransfer = "user-transfer"
	ErrorMessage = "error-message"

)

type Action struct {
	Data interface{} `json:"data"`
	Time int64	`json:"time"`
	Action string `json:"action"`
}

func (action *Action) Marshal() (b []byte, err error) {
	if action.Action == PingAction {
		return []byte(""), nil
	}
	if action.Action == ReceiveMessageAction {
		msg, ok := action.Data.(*models.Message)
		if !ok {
			err = errors.New("param error")
			return
		}
		b, err = json.Marshal(Action{
			Time:   action.Time,
			Action: action.Action,
			Data: msg.ToJson(),
		})
		return
	}
	b, err = json.Marshal(action)
	return
}

func (action *Action) UnMarshal(b []byte) (err error) {
	err =  json.Unmarshal(b, action)
	return
}
// 获取action的message
func (action *Action) GetMessage() (message *models.Message,err error)  {
	if action.Action == SendMessageAction {
		message = &models.Message{}
		err = mapstructure.Decode(action.Data, message)
	} else {
		err = errors.New("invalid action")
	}
	return

}
func NewReceiveAction(msg *models.Message) *Action {
	return &Action{
		Action: ReceiveMessageAction,
		Time: time.Now().Unix(),
		Data: msg,
	}
}
func NewReceiptAction(msg *models.Message) (act *Action) {
	data := make(map[string]interface{})
	data["user_id"] = msg.UserId
	data["req_id"] = msg.ReqId
	act = &Action{
		Action: ReceiptAction,
		Time: time.Now().Unix(),
		Data: data,
	}
	return
}
func NewAdminsAction(admins []resource.Admin) *Action {
	return &Action{
		Action: AdminsAction,
		Time: time.Now().Unix(),
		Data: admins,
	}
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
func NewMoreThanOne() *Action {
	return &Action{
		Action: MoreThanOne,
		Time: time.Now().Unix(),
	}
}
func NewOtherLogin() *Action {
	return &Action{
		Action: OtherLogin,
		Time: time.Now().Unix(),
	}
}
func NewPing() *Action {
	return &Action{
		Action: PingAction,
		Time: time.Now().Unix(),
	}
}
func NewWaitingUsers(i interface{}) *Action {
	return &Action{
		Action: WaitingUserAction,
		Time: time.Now().Unix(),
		Data: i,
	}
}
func NewWaitingUserCount(count int64) *Action {
	return &Action{
		Data:   count,
		Time:   time.Now().Unix(),
		Action: WaitingUserCount,
	}
}
func NewUserTransfer(i interface{}) *Action {
	return &Action{
		Data:   i,
		Time:   time.Now().Unix(),
		Action: UserTransfer,
	}
}
func NewErrorMessage(msg string)  *Action {
	return &Action{
		Data:   msg,
		Time:   time.Now().Unix(),
		Action: ErrorMessage,
	}
}

