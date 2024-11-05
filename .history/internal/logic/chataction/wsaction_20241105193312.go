package Action

import (
	"encoding/json"
	"errors"
	"gf-chat/internal/consts"
	"gf-chat/internal/model/chat"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/model/relation"
	"gf-chat/internal/service"
	"time"

	"github.com/gogf/gf/v2/util/gconv"
)

func init() {
	service.RegisterAction(&sAction{})
}

type sAction struct {
}

func (s sAction) GetMessage(action *chat.Action) (message *relation.CustomerChatMessages, err error) {
	if action.Action == consts.ActionSendMessage {
		message = &relation.CustomerChatMessages{}
		err = gconv.Struct(action.Data, message)
	} else {
		err = errors.New("invalid action")
	}
	return
}

func (s sAction) UnMarshalAction(b []byte) (action *chat.Action, err error) {
	action = &chat.Action{}
	err = json.Unmarshal(b, action)
	return
}
func (s sAction) MarshalAction(action *chat.Action) (b []byte, err error) {
	if action.Action == consts.ActionPing {
		return []byte(""), nil
	}
	if action.Action == consts.ActionReceiveMessage {
		msg, ok := action.Data.(*relation.CustomerChatMessages)
		if !ok {
			err = errors.New("param error")
			return
		}
		action.Data = service.ChatMessage().RelationToChat(*msg)
	}
	b, err = json.Marshal(action)
	return
}

func (s sAction) String(action *chat.Action) string {
	b, err := json.Marshal(action)
	if err == nil {
		return string(b)
	}
	return ""
}

func (s sAction) NewReceiveAction(msg *relation.CustomerChatMessages) *chat.Action {
	return &chat.Action{
		Action: consts.ActionReceiveMessage,
		Time:   time.Now().Unix(),
		Data:   msg,
	}
}
func (s sAction) NewReceiptAction(msg *relation.CustomerChatMessages) (act *chat.Action) {
	data := make(map[string]interface{})
	data["user_id"] = msg.UserId
	data["req_id"] = msg.ReqId
	data["msg_id"] = msg.Id
	act = &chat.Action{
		Action: consts.ActionReceipt,
		Time:   time.Now().Unix(),
		Data:   data,
	}
	return
}
func (s sAction) NewAdminsAction(d any) (act *chat.Action) {
	return &chat.Action{
		Action: consts.ActionAdmins,
		Time:   time.Now().Unix(),
		Data:   d,
	}
}
func (s sAction) NewUserOnline(uid int, platform string) *chat.Action {
	data := make(map[string]interface{})
	data["user_id"] = uid
	data["platform"] = platform
	return &chat.Action{
		Action: consts.ActionUserOnLine,
		Time:   time.Now().Unix(),
		Data:   data,
	}
}
func (s sAction) NewUserOffline(uid uint) *chat.Action {
	data := make(map[string]interface{})
	data["user_id"] = uid
	return &chat.Action{
		Action: consts.ActionUserOffLine,
		Time:   time.Now().Unix(),
		Data:   data,
	}
}
func (s sAction) NewMoreThanOne() *chat.Action {
	return &chat.Action{
		Action: consts.ActionMoreThanOne,
		Time:   time.Now().Unix(),
	}
}
func (s sAction) NewOtherLogin() *chat.Action {
	return &chat.Action{
		Action: consts.ActionOtherLogin,
		Time:   time.Now().Unix(),
	}
}
func (s sAction) NewPing() *chat.Action {
	return &chat.Action{
		Action: consts.ActionPing,
		Time:   time.Now().Unix(),
	}
}
func (s sAction) NewWaitingUsers(i interface{}) *chat.Action {
	return &chat.Action{
		Action: consts.ActionWaitingUser,
		Time:   time.Now().Unix(),
		Data:   i,
	}
}
func (s sAction) NewWaitingUserCount(count uint) *chat.Action {
	return &chat.Action{
		Data:   count,
		Time:   time.Now().Unix(),
		Action: consts.ActionWaitingUserCount,
	}
}
func (s sAction) NewUserTransfer(i interface{}) *chat.Action {
	return &chat.Action{
		Data:   i,
		Time:   time.Now().Unix(),
		Action: consts.ActionUserTransfer,
	}
}
func (s sAction) NewErrorMessage(msg string) *chat.Action {
	return &chat.Action{
		Data:   msg,
		Time:   time.Now().Unix(),
		Action: consts.ActionErrorMessage,
	}
}

func (s sAction) NewReadAction(msgIds []int64) *chat.Action {
	return &chat.Action{
		Data:   msgIds,
		Time:   time.Now().Unix(),
		Action: consts.ActionRead,
	}
}
func (s sAction) NewRateAction(message *entity.CustomerChatMessages) *chat.Action {
	data := make(map[string]interface{})
	data["msg_id"] = message.Id
	data["rate"] = message.Content
	data["user_id"] = message.UserId
	return &chat.Action{
		Action: consts.ActionUserRate,
		Time:   time.Now().Unix(),
		Data:   data,
	}
}
