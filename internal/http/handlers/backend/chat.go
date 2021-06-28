package backend

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
	"ws/internal/auth"
	"ws/internal/chat"
	"ws/internal/file"
	"ws/internal/json"
	"ws/internal/repositories"
	"ws/internal/util"
	"ws/internal/websocket"
)

// 获取消息
func GetHistoryMessage(c *gin.Context) {
	var uid int64
	var mid int64
	var err error
	uidStr, exist := c.GetQuery("uid")
	if !exist {
		util.RespValidateFail(c, "invalid params")
		return
	}
	uid, err = strconv.ParseInt(uidStr, 10, 64)
	if err != nil {
		util.RespValidateFail(c, "invalid params")
		return
	}
	user, exist := repositories.GetUserById(uid)
	if !exist {
		util.RespValidateFail(c, "invalid params")
		return
	}
	backendUser := auth.GetBackendUser(c)
	wheres := []*repositories.Where{
		{
			Filed: "service_id = ?",
			Value: backendUser.GetPrimaryKey(),
		},
		{
			Filed: "user_id = ?",
			Value: uid,
		},
	}
	midStr, exist := c.GetQuery("mid")
	if exist {
		mid, err = strconv.ParseInt(midStr, 10, 64)
		if err == nil {
			wheres = append(wheres, &repositories.Where{
				Filed: "id < ?",
				Value: mid,
			})
		}
	}
	messages := repositories.GetMessages(wheres, 20, []string{"User"})
	res := make([]*json.Message, 0)
	for _, m := range messages {
		var avatar string
		if m.Source == 1 {
			avatar = backendUser.GetAvatarUrl()
		} else if m.Source == 0 {
			avatar = user.GetAvatarUrl()
		} else if m.Source == 2 {
			avatar = util.PublicAsset("avatar.jpeg")
		}
		res = append(res, &json.Message{
			Id:         m.Id,
			UserId:     m.UserId,
			ServiceId:  m.ServiceId,
			Type:       m.Type,
			Content:    m.Content,
			IsSuccess:  true,
			ReceivedAT: m.ReceivedAT,
			Source:   m.Source,
			ReqId:      m.ReqId,
			IsRead:     m.IsRead,
			Avatar:    avatar,
		})
	}
	util.RespSuccess(c, res)
}

// 聊天用户列表
func ChatUserList(c *gin.Context) {
	backendUser := auth.GetBackendUser(c)
	ids, times := chat.GetChatUserIds(backendUser.GetPrimaryKey())
	users := repositories.GetUserByIds(ids)
	resp := make([]*json.User, 0, len(users))
	userMap := make(map[int64]auth.User)
	for _, user := range users {
		userMap[user.GetPrimaryKey()] = user
	}
	deadline := chat.GetDeadlineTime()
	for index, id := range ids {
		u := userMap[id]
		chatUserRes := &json.User{
			ID:       u.GetPrimaryKey(),
			Username: u.GetUsername(),
			Messages: make([]*json.Message, 0),
			Unread:   0,
		}
		lastChatTime := times[index]
		chatUserRes.LastChatTime =lastChatTime
		chatUserRes.Disabled = lastChatTime < deadline
		if _, ok := websocket.UserHub.GetConn(u.GetPrimaryKey()); ok {
			chatUserRes.Online = true

		}
		resp = append(resp, chatUserRes)
	}
	messages := repositories.GetMessages([]*repositories.Where{
		{
			Filed: "received_at > ?",
			Value: deadline,
		},
		{
			Filed: "service_id = ?",
			Value: backendUser.GetPrimaryKey(),
		},
	}, -1, []string{"User", "BackendUser"})
	for _, u := range resp {
		for _, m := range messages {
			if m.UserId == u.ID {
				var avatar string
				if m.Source == 0 {
					avatar = userMap[u.ID].GetAvatarUrl()
				} else if m.Source == 1 {
					avatar = backendUser.GetAvatarUrl()
				}else if m.Source == 2 {
					avatar = util.PublicAsset("avatar.jpeg")
				}
				rm := &json.Message{
					Id:         m.Id,
					UserId:     m.UserId,
					ServiceId:  m.ServiceId,
					Type:       m.Type,
					Content:    m.Content,
					ReceivedAT: m.ReceivedAT,
					Source:   m.Source,
					ReqId:      m.ReqId,
					IsRead:     m.IsRead,
					Avatar:     avatar,
				}
				rm.IsSuccess = true
				if !m.IsRead && m.Source == 0 {
					u.Unread += 1
				}
				u.Messages = append(u.Messages, rm)
			}
		}
		if _, ok := websocket.UserHub.GetConn(u.ID); ok {
			u.Online = true
		}
	}
	util.RespSuccess(c, resp)
}

// 接入用户
func AcceptUser(c *gin.Context) {
	form := &struct {
		Uid int64
	}{}
	err := c.Bind(form)
	if err != nil {
		util.RespValidateFail(c, "invalid params")
		return
	}
	user, exist := repositories.GetUserById(form.Uid)
	if !exist {
		util.RespValidateFail(c, "invalid params")
		return
	}
	if chat.GetUserLastServerId(user.GetPrimaryKey()) != 0 {
		util.RespFail(c, "frontend had been accepted", 10001)
		return
	}
	unSendMsg := repositories.GetUnSendMessage([]*repositories.Where{
		{
			Filed: "user_id = ?",
			Value: user.GetPrimaryKey(),
		},
	}, []string{"User"})
	backendUser := auth.GetBackendUser(c)
	now := time.Now().Unix()
	// 更新未发送的消息
	repositories.UpdateMessages([]*repositories.Where{
		{
			Filed: "user_id = ?",
			Value: user.GetPrimaryKey(),
		},
		{
			Filed: "service_id = ?",
			Value: 0,
		},
	}, map[string]interface{}{
		"service_id": backendUser.ID,
		"send_at":    now,
	})
	messages := repositories.GetMessages([]*repositories.Where{
		{
			Filed: "user_id = ?",
			Value: user.GetPrimaryKey(),
		},
		{
			Filed: "service_id = ?",
			Value: backendUser.GetPrimaryKey(),
		},
	}, 20, []string{})
	_ = chat.SetUserServerId(user.GetPrimaryKey(), backendUser.GetPrimaryKey())
	messageLength := len(messages)
	chatUser := &json.User{
		ID:           user.GetPrimaryKey(),
		Username:     user.GetUsername(),
		LastChatTime: 0,
		Messages:     make([]*json.Message, messageLength, messageLength),
	}
	chatUser.Unread = len(unSendMsg)
	_, exist = websocket.UserHub.GetConn(user.GetPrimaryKey())
	chatUser.Online = exist
	chatUser.LastChatTime = time.Now().Unix()
	for index, m := range messages {
		rm := &json.Message{
			Id:         m.Id,
			UserId:     m.UserId,
			ServiceId:  m.ServiceId,
			Type:       m.Type,
			Content:    m.Content,
			ReceivedAT: m.ReceivedAT,
			Source:   m.Source,
			ReqId:      m.ReqId,
			IsSuccess:  true,
			IsRead:     m.IsRead,
			Avatar:     user.GetAvatarUrl(),
		}
		chatUser.Messages[index] = rm
	}
	go websocket.ServiceHub.BroadcastWaitingUser()
	util.RespSuccess(c, chatUser)
}

// 移除用户
func RemoveUser(c *gin.Context) {
	uidStr := c.Param("id")
	uid, err := strconv.ParseInt(uidStr, 10, 64)
	if err == nil {
		BackendUser := auth.GetBackendUser(c)
		_ = chat.RemoveUserServerId(uid, BackendUser.GetPrimaryKey())
	}
	util.RespSuccess(c, nil)
}

// 已读
func ReadAll(c *gin.Context) {
	form := &struct {
		Id int64
	}{}
	err := c.Bind(form)
	if err == nil {
		backendUser := auth.GetBackendUser(c)
		wheres := []*repositories.Where{
			{
				Filed: "service_id = ?",
				Value: backendUser.GetPrimaryKey(),
			},
			{
				Filed: "user_id = ?",
				Value: form.Id,
			},
			{
				Filed: "is_read = ?",
				Value: 0,
			},
		}
		repositories.UpdateMessages(wheres, map[string]interface{}{
			"is_read": 1,
		})
		util.RespSuccess(c, gin.H{})
	} else {
		util.RespValidateFail(c, "invalid params")
	}
}

// 聊天图片
func Image(c *gin.Context) {
	f, _ := c.FormFile("file")
	ff, err := file.Save(f, "chat")
	if err != nil {
		util.RespFail(c, err.Error(), 500)
	} else {
		util.RespSuccess(c, gin.H{
			"url": ff.FullUrl,
		})
	}
}
