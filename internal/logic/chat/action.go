package chat

import (
	"context"
	"encoding/json"
	"gf-chat/api"
	"gf-chat/internal/consts"
	"gf-chat/internal/model"
	"gf-chat/internal/service"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"github.com/gogf/gf/v2/util/gconv"
)

var action = &iAction{}

type iAction struct{}

// getMessage extracts message from action
func (a *iAction) getMessage(action *api.ChatAction) (*model.CustomerChatMessage, error) {
	if action.Action != consts.ActionSendMessage {
		return nil, gerror.NewCode(gcode.CodeValidationFailed, "invalid action")
	}

	message := &model.CustomerChatMessage{}
	if err := gconv.Struct(action.Data, message); err != nil {
		return nil, err
	}
	return message, nil
}

// unMarshal unmarshals JSON bytes into ChatAction
func (a *iAction) unMarshal(b []byte) (*api.ChatAction, error) {
	action := &api.ChatAction{}
	if err := json.Unmarshal(b, action); err != nil {
		return nil, err
	}
	return action, nil
}

// marshal marshals ChatAction to JSON bytes
func (a *iAction) marshal(ctx context.Context, action api.ChatAction) ([]byte, error) {
	if action.Action == consts.ActionPing {
		return []byte("ping"), nil
	}

	if action.Action == consts.ActionReceiveMessage {
		msg, ok := action.Data.(*model.CustomerChatMessage)
		if !ok {
			return nil, gerror.NewCode(gcode.CodeValidationFailed, "param error")
		}
		data, err := service.ChatMessage().ToApi(ctx, msg)
		if err != nil {
			return nil, err
		}
		action.Data = data
	}

	return json.Marshal(action)
}

// newAction creates a new ChatAction with common fields
func (a *iAction) newAction(actionType string, data interface{}) *api.ChatAction {
	return &api.ChatAction{
		Action: actionType,
		Time:   time.Now().Unix(),
		Data:   data,
	}
}

// newReceive creates a new receive message action
func (a *iAction) newReceive(msg *model.CustomerChatMessage) *api.ChatAction {
	return a.newAction(consts.ActionReceiveMessage, msg)
}

// newReceipt creates a new receipt action
func (a *iAction) newReceipt(msg *model.CustomerChatMessage) *api.ChatAction {
	data := map[string]interface{}{
		"user_id": msg.UserId,
		"req_id":  msg.ReqId,
		"msg_id":  msg.Id,
	}
	return a.newAction(consts.ActionReceipt, data)
}

// newAdmins creates a new admins action
func (a *iAction) newAdmins(d interface{}) *api.ChatAction {
	return a.newAction(consts.ActionAdmins, d)
}

// newUserOnline creates a new user online action
func (a *iAction) newUserOnline(uid uint, platform string) *api.ChatAction {
	data := map[string]interface{}{
		"user_id":  uid,
		"platform": platform,
	}
	return a.newAction(consts.ActionUserOnLine, data)
}

// newUserOffline creates a new user offline action
func (a *iAction) newUserOffline(uid uint) *api.ChatAction {
	return a.newAction(consts.ActionUserOffLine, map[string]interface{}{"user_id": uid})
}

// newMoreThanOne creates a new more than one action
func (a *iAction) newMoreThanOne() *api.ChatAction {
	return a.newAction(consts.ActionMoreThanOne, nil)
}

// newOtherLogin creates a new other login action
func (a *iAction) newOtherLogin() *api.ChatAction {
	return a.newAction(consts.ActionOtherLogin, nil)
}

// newPing creates a new ping action
func (a *iAction) newPing() *api.ChatAction {
	return a.newAction(consts.ActionPing, nil)
}

// newWaitingUsers creates a new waiting users action
func (a *iAction) newWaitingUsers(i interface{}) *api.ChatAction {
	return a.newAction(consts.ActionWaitingUser, i)
}

// newWaitingUserCount creates a new waiting user count action
func (a *iAction) newWaitingUserCount(count uint) *api.ChatAction {
	return a.newAction(consts.ActionWaitingUserCount, count)
}

// newUserTransfer creates a new user transfer action
func (a *iAction) newUserTransfer(i interface{}) *api.ChatAction {
	return a.newAction(consts.ActionUserTransfer, i)
}

// newErrorMessage creates a new error message action
func (a *iAction) newErrorMessage(msg string) *api.ChatAction {
	return a.newAction(consts.ActionErrorMessage, msg)
}

// newReadAction creates a new read action
func (a *iAction) newReadAction(msgIds []uint) *api.ChatAction {
	return a.newAction(consts.ActionRead, msgIds)
}

// newRateAction creates a new rate action
func (a *iAction) newRateAction(message *model.CustomerChatMessage) *api.ChatAction {
	data := map[string]interface{}{
		"msg_id":  message.Id,
		"rate":    message.Content,
		"user_id": message.UserId,
	}
	return a.newAction(consts.ActionUserRate, data)
}
