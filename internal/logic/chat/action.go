package chat

import (
	"encoding/json"
	"errors"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/consts"
	"gf-chat/internal/model"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/os/gctx"
	"time"

	"github.com/gogf/gf/v2/util/gconv"
)

func GetMessage(action *api.ChatAction) (message *model.CustomerChatMessage, err error) {
	if action.Action == consts.ActionSendMessage {
		message = &model.CustomerChatMessage{}
		err = gconv.Struct(action.Data, message)
	} else {
		err = errors.New("invalid action")
	}
	return
}

func unMarshalAction(b []byte) (action *api.ChatAction, err error) {
	action = &api.ChatAction{}
	err = json.Unmarshal(b, action)
	return
}
func marshalAction(action *api.ChatAction) (b []byte, err error) {
	if action.Action == consts.ActionPing {
		return []byte(""), nil
	}
	if action.Action == consts.ActionReceiveMessage {
		msg, ok := action.Data.(*model.CustomerChatMessage)
		if !ok {
			err = errors.New("param error")
			return
		}
		data, err := service.ChatMessage().RelationToChat(gctx.New(), *msg)
		if err != nil {
			return b, err
		}
		action.Data = data
	}
	b, err = json.Marshal(action)
	return
}

func newReceiveAction(msg *model.CustomerChatMessage) *api.ChatAction {
	return &api.ChatAction{
		Action: consts.ActionReceiveMessage,
		Time:   time.Now().Unix(),
		Data:   msg,
	}
}
func newReceiptAction(msg *model.CustomerChatMessage) (act *api.ChatAction) {
	data := make(map[string]interface{})
	data["user_id"] = msg.UserId
	data["req_id"] = msg.ReqId
	data["msg_id"] = msg.Id
	act = &api.ChatAction{
		Action: consts.ActionReceipt,
		Time:   time.Now().Unix(),
		Data:   data,
	}
	return
}
func newAdminsAction(d any) (act *api.ChatAction) {
	return &api.ChatAction{
		Action: consts.ActionAdmins,
		Time:   time.Now().Unix(),
		Data:   d,
	}
}
func newUserOnlineAction(uid uint, platform string) *api.ChatAction {
	data := make(map[string]interface{})
	data["user_id"] = uid
	data["platform"] = platform
	return &api.ChatAction{
		Action: consts.ActionUserOnLine,
		Time:   time.Now().Unix(),
		Data:   data,
	}
}
func newUserOfflineAction(uid uint) *api.ChatAction {
	data := make(map[string]interface{})
	data["user_id"] = uid
	return &api.ChatAction{
		Action: consts.ActionUserOffLine,
		Time:   time.Now().Unix(),
		Data:   data,
	}
}
func newMoreThanOneAction() *api.ChatAction {
	return &api.ChatAction{
		Action: consts.ActionMoreThanOne,
		Time:   time.Now().Unix(),
	}
}
func newOtherLoginAction() *api.ChatAction {
	return &api.ChatAction{
		Action: consts.ActionOtherLogin,
		Time:   time.Now().Unix(),
	}
}
func newPingAction() *api.ChatAction {
	return &api.ChatAction{
		Action: consts.ActionPing,
		Time:   time.Now().Unix(),
	}
}
func newWaitingUsersAction(i interface{}) *api.ChatAction {
	return &api.ChatAction{
		Action: consts.ActionWaitingUser,
		Time:   time.Now().Unix(),
		Data:   i,
	}
}
func newWaitingUserCountAction(count uint) *api.ChatAction {
	return &api.ChatAction{
		Data:   count,
		Time:   time.Now().Unix(),
		Action: consts.ActionWaitingUserCount,
	}
}
func newUserTransferAction(i interface{}) *api.ChatAction {
	return &api.ChatAction{
		Data:   i,
		Time:   time.Now().Unix(),
		Action: consts.ActionUserTransfer,
	}
}
func newErrorMessageAction(msg string) *api.ChatAction {
	return &api.ChatAction{
		Data:   msg,
		Time:   time.Now().Unix(),
		Action: consts.ActionErrorMessage,
	}
}

func newReadActionAction(msgIds []uint) *api.ChatAction {
	return &api.ChatAction{
		Data:   msgIds,
		Time:   time.Now().Unix(),
		Action: consts.ActionRead,
	}
}
func newRateActionAction(message *model.CustomerChatMessage) *api.ChatAction {
	data := make(map[string]interface{})
	data["msg_id"] = message.Id
	data["rate"] = message.Content
	data["user_id"] = message.UserId
	return &api.ChatAction{
		Action: consts.ActionUserRate,
		Time:   time.Now().Unix(),
		Data:   data,
	}
}
