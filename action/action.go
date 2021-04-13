package action

import (
	"encoding/json"
	"errors"
	"github.com/mitchellh/mapstructure"
	"time"
	"ws/models"
	"ws/util"
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
	ReqId int64 `json:"req_id"`
	Data map[string]interface{} `json:"data"`
	Time int64	`json:"time"`
	Action string `json:"action"`
	Message *models.Message `json:"-"`
}

func (action *Action) Marshal() (b []byte, err error) {
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
func NewMessage(message models.Message) *Action {
	d := make(map[string]interface{})
	d["id"] = message.Id
	d["user_id"] = message.UserId
	d["type"] = message.Type
	d["content"] = message.Content
	d["is_server"] = message.IsServer
	return &Action{
		ReqId: util.CreateReqId(),
		Action: MessageAction,
		Time: time.Now().Unix(),
		Data: d,
	}
}
func NewServerUserList(d []*models.ChatUser) *Action {
	data := make(map[string]interface{})
	data["list"] = d
	return &Action{
		ReqId: util.CreateReqId(),
		Action: ServerUserListAction,
		Time: time.Now().Unix(),
		Data: data,
	}
}
func NewReceipt(action *Action) *Action {
	data := make(map[string]interface{})
	userId, ok := action.Data["user_id"]
	if ok {
		data["user_id"] = userId
	}
	return &Action{
		ReqId: action.ReqId,
		Action: ReceiptAction,
		Time: time.Now().Unix(),
		Data: data,
	}
}
func NewUserOnline(uid int64) *Action {
	data := make(map[string]interface{})
	data["user_id"] = uid
	return &Action{
		ReqId: util.CreateReqId(),
		Action: UserOnLineAction,
		Time: time.Now().Unix(),
		Data: data,
	}
}
func NewUserOffline(uid int64) *Action {
	data := make(map[string]interface{})
	data["user_id"] = uid
	return &Action{
		ReqId: util.CreateReqId(),
		Action: UserOffLineAction,
		Time: time.Now().Unix(),
		Data: data,
	}
}
func NewPing() *Action {
	return &Action{
		ReqId: util.CreateReqId(),
		Action: PingAction,
		Time: time.Now().Unix(),
	}
}
func NewWaitingUsers(i interface{}) *Action {
	data := make(map[string]interface{})
	data["list"] = i
	return &Action{
		ReqId: util.CreateReqId(),
		Action: UserWaitingCountAction,
		Time: time.Now().Unix(),
		Data: data,
	}
}
func NewServiceOnlineList(data map[string]interface{}) *Action {
	return &Action{
		ReqId: util.CreateReqId(),
		Action: "service_online_list",
		Time: time.Now().Unix(),
		Data: data,
	}
}

