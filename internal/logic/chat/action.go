package chat

import (
	"context"
	"encoding/json"
	"gf-chat/api"
	"gf-chat/internal/consts"
	"gf-chat/internal/model"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/errors/gerror"
	"time"

	"github.com/gogf/gf/v2/util/gconv"
)

var action = &iaction{}

type iaction struct {
}

func (a iaction) getMessage(action *api.ChatAction) (message *model.CustomerChatMessage, err error) {
	if action.Action == consts.ActionSendMessage {
		message = &model.CustomerChatMessage{}
		err = gconv.Struct(action.Data, message)
	} else {
		err = gerror.New("invalid action")
	}
	return
}

func (a iaction) unMarshal(b []byte) (action *api.ChatAction, err error) {
	action = &api.ChatAction{}
	err = json.Unmarshal(b, action)
	return
}
func (a iaction) marshal(ctx context.Context, action api.ChatAction) (b []byte, err error) {
	if action.Action == consts.ActionPing {
		return []byte(""), nil
	}
	if action.Action == consts.ActionReceiveMessage {
		msg, ok := action.Data.(*model.CustomerChatMessage)
		if !ok {
			err = gerror.New("param error")
			return
		}
		data, err := service.ChatMessage().ToApi(ctx, msg)
		if err != nil {
			return b, err
		}
		action.Data = data
	}
	b, err = json.Marshal(action)
	return
}

func (a iaction) newReceive(msg *model.CustomerChatMessage) *api.ChatAction {
	return &api.ChatAction{
		Action: consts.ActionReceiveMessage,
		Time:   time.Now().Unix(),
		Data:   msg,
	}
}
func (a iaction) newReceipt(msg *model.CustomerChatMessage) (act *api.ChatAction) {
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
func (a iaction) newAdmins(d any) (act *api.ChatAction) {
	return &api.ChatAction{
		Action: consts.ActionAdmins,
		Time:   time.Now().Unix(),
		Data:   d,
	}
}
func (a iaction) newUserOnline(uid uint, platform string) *api.ChatAction {
	data := make(map[string]interface{})
	data["user_id"] = uid
	data["platform"] = platform
	return &api.ChatAction{
		Action: consts.ActionUserOnLine,
		Time:   time.Now().Unix(),
		Data:   data,
	}
}
func (a iaction) newUserOffline(uid uint) *api.ChatAction {
	data := make(map[string]interface{})
	data["user_id"] = uid
	return &api.ChatAction{
		Action: consts.ActionUserOffLine,
		Time:   time.Now().Unix(),
		Data:   data,
	}
}
func (a iaction) newMoreThanOne() *api.ChatAction {
	return &api.ChatAction{
		Action: consts.ActionMoreThanOne,
		Time:   time.Now().Unix(),
	}
}
func (a iaction) newOtherLogin() *api.ChatAction {
	return &api.ChatAction{
		Action: consts.ActionOtherLogin,
		Time:   time.Now().Unix(),
	}
}
func (a iaction) newPing() *api.ChatAction {
	return &api.ChatAction{
		Action: consts.ActionPing,
		Time:   time.Now().Unix(),
	}
}
func (a iaction) newWaitingUsers(i interface{}) *api.ChatAction {
	return &api.ChatAction{
		Action: consts.ActionWaitingUser,
		Time:   time.Now().Unix(),
		Data:   i,
	}
}
func (a iaction) newWaitingUserCount(count uint) *api.ChatAction {
	return &api.ChatAction{
		Data:   count,
		Time:   time.Now().Unix(),
		Action: consts.ActionWaitingUserCount,
	}
}
func (a iaction) newUserTransfer(i interface{}) *api.ChatAction {
	return &api.ChatAction{
		Data:   i,
		Time:   time.Now().Unix(),
		Action: consts.ActionUserTransfer,
	}
}
func (a iaction) newErrorMessage(msg string) *api.ChatAction {
	return &api.ChatAction{
		Data:   msg,
		Time:   time.Now().Unix(),
		Action: consts.ActionErrorMessage,
	}
}

func (a iaction) newReadAction(msgIds []uint) *api.ChatAction {
	return &api.ChatAction{
		Data:   msgIds,
		Time:   time.Now().Unix(),
		Action: consts.ActionRead,
	}
}
func (a iaction) newRateAction(message *model.CustomerChatMessage) *api.ChatAction {
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
