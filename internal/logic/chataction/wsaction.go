package Action

import (
	"encoding/json"
	"errors"
	"gf-chat/internal/consts"
	"gf-chat/internal/model"
	"gf-chat/internal/service"
	"time"

	"github.com/gogf/gf/v2/util/gconv"
)

func init() {
	service.RegisterAction(&sAction{})
}

type sAction struct {
}

func (s sAction) GetMessage(action *model.ChatAction) (message *model.CustomerChatMessage, err error) {
	if action.Action == consts.ActionSendMessage {
		message = &model.CustomerChatMessage{}
		err = gconv.Struct(action.Data, message)
	} else {
		err = errors.New("invalid action")
	}
	return
}

func (s sAction) UnMarshalAction(b []byte) (action *model.ChatAction, err error) {
	action = &model.ChatAction{}
	err = json.Unmarshal(b, action)
	return
}
func (s sAction) MarshalAction(action *model.ChatAction) (b []byte, err error) {
	if action.Action == consts.ActionPing {
		return []byte(""), nil
	}
	if action.Action == consts.ActionReceiveMessage {
		msg, ok := action.Data.(*model.CustomerChatMessage)
		if !ok {
			err = errors.New("param error")
			return
		}
		action.Data = service.ChatMessage().RelationToChat(*msg)
	}
	b, err = json.Marshal(action)
	return
}

func (s sAction) String(action *model.ChatAction) string {
	b, err := json.Marshal(action)
	if err == nil {
		return string(b)
	}
	return ""
}

func (s sAction) NewReceiveAction(msg *model.CustomerChatMessage) *model.ChatAction {
	return &model.ChatAction{
		Action: consts.ActionReceiveMessage,
		Time:   time.Now().Unix(),
		Data:   msg,
	}
}
func (s sAction) NewReceiptAction(msg *model.CustomerChatMessage) (act *model.ChatAction) {
	data := make(map[string]interface{})
	data["user_id"] = msg.UserId
	data["req_id"] = msg.ReqId
	data["msg_id"] = msg.Id
	act = &model.ChatAction{
		Action: consts.ActionReceipt,
		Time:   time.Now().Unix(),
		Data:   data,
	}
	return
}
func (s sAction) NewAdminsAction(d any) (act *model.ChatAction) {
	return &model.ChatAction{
		Action: consts.ActionAdmins,
		Time:   time.Now().Unix(),
		Data:   d,
	}
}
func (s sAction) NewUserOnline(uid uint, platform string) *model.ChatAction {
	data := make(map[string]interface{})
	data["user_id"] = uid
	data["platform"] = platform
	return &model.ChatAction{
		Action: consts.ActionUserOnLine,
		Time:   time.Now().Unix(),
		Data:   data,
	}
}
func (s sAction) NewUserOffline(uid uint) *model.ChatAction {
	data := make(map[string]interface{})
	data["user_id"] = uid
	return &model.ChatAction{
		Action: consts.ActionUserOffLine,
		Time:   time.Now().Unix(),
		Data:   data,
	}
}
func (s sAction) NewMoreThanOne() *model.ChatAction {
	return &model.ChatAction{
		Action: consts.ActionMoreThanOne,
		Time:   time.Now().Unix(),
	}
}
func (s sAction) NewOtherLogin() *model.ChatAction {
	return &model.ChatAction{
		Action: consts.ActionOtherLogin,
		Time:   time.Now().Unix(),
	}
}
func (s sAction) NewPing() *model.ChatAction {
	return &model.ChatAction{
		Action: consts.ActionPing,
		Time:   time.Now().Unix(),
	}
}
func (s sAction) NewWaitingUsers(i interface{}) *model.ChatAction {
	return &model.ChatAction{
		Action: consts.ActionWaitingUser,
		Time:   time.Now().Unix(),
		Data:   i,
	}
}
func (s sAction) NewWaitingUserCount(count uint) *model.ChatAction {
	return &model.ChatAction{
		Data:   count,
		Time:   time.Now().Unix(),
		Action: consts.ActionWaitingUserCount,
	}
}
func (s sAction) NewUserTransfer(i interface{}) *model.ChatAction {
	return &model.ChatAction{
		Data:   i,
		Time:   time.Now().Unix(),
		Action: consts.ActionUserTransfer,
	}
}
func (s sAction) NewErrorMessage(msg string) *model.ChatAction {
	return &model.ChatAction{
		Data:   msg,
		Time:   time.Now().Unix(),
		Action: consts.ActionErrorMessage,
	}
}

func (s sAction) NewReadAction(msgIds []uint) *model.ChatAction {
	return &model.ChatAction{
		Data:   msgIds,
		Time:   time.Now().Unix(),
		Action: consts.ActionRead,
	}
}
func (s sAction) NewRateAction(message *model.CustomerChatMessage) *model.ChatAction {
	data := make(map[string]interface{})
	data["msg_id"] = message.Id
	data["rate"] = message.Content
	data["user_id"] = message.UserId
	return &model.ChatAction{
		Action: consts.ActionUserRate,
		Time:   time.Now().Unix(),
		Data:   data,
	}
}
