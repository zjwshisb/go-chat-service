package admin

import (
	"github.com/gin-gonic/gin"
	"sort"
	"strconv"
	"time"
	"ws/app/auth"
	"ws/app/chat"
	"ws/app/file"
	"ws/app/json"
	"ws/app/models"
	"ws/app/repositories"
	"ws/app/util"
	"ws/app/websocket"
)

type ChatHandler struct {
}

// 查看历史对话
func (handle *ChatHandler) GetHistorySession(c *gin.Context) {
	useId := c.Param("uid")
	sessions := chatSessionRepo.Get([]Where{
		{
			Filed: "user_id = ?",
			Value: useId,
		},
		{
			Filed: "admin_id > ?",
			Value: 0,
		},
	}, -1, []string{"Admin","User"}, "id desc")
	resp := make([]*json.ChatSession, len(sessions))
	for i, session := range sessions {
		resp[i] = session.ToJson()
	}
	util.RespSuccess(c, resp)
}

// 获取消息
func (handle *ChatHandler) GetHistoryMessage(c *gin.Context) {
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
	admin := auth.GetAdmin(c)
	chatIds, _ := chat.AdminService.GetUsersWithLimitTime(admin.GetPrimaryKey())
	userExist := false
	for _, chatId := range  chatIds {
		if chatId == uid {
			userExist = true
		}
	}
	if !userExist {
		util.RespValidateFail(c, "invalid params")
		return
	}
	wheres := []Where{
		{
			Filed: "admin_id = ?",
			Value: admin.GetPrimaryKey(),
		},
		{
			Filed: "user_id = ?",
			Value: uid,
		},
		{
			Filed: "source in ?",
			Value: []int{models.SourceAdmin, models.SourceUser},
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
	messages := messageRepo.Get(wheres, 20, []string{"User", "Admin"}, "id desc")
	res := make([]*models.MessageJson, 0)
	for _, m := range messages {
		res = append(res, m.ToJson())
	}
	util.RespSuccess(c, res)
}

// 聊天用户列表
func (handle *ChatHandler) ChatUserList(c *gin.Context) {
	admin := auth.GetAdmin(c)
	ids, times := chat.AdminService.GetUsersWithLimitTime(admin.GetPrimaryKey())
	users := userRepo.Get([]Where{
		{
			Filed: "id in ?",
			Value: ids,
		},
	}, -1, []string{})
	resp := make([]*models.UserJson, 0, len(users))
	userMap := make(map[int64]auth.User)
	for _, user := range users {
		userMap[user.GetPrimaryKey()] = user
	}
	for index, id := range ids {
		limitTime := times[index]
		disabled := limitTime <= time.Now().Unix()
		// 聊天列表超过50时，不显示已失效的用户
		if len(resp) >= 50 && disabled {
			go func() {
				_ = chat.AdminService.RemoveUser(admin.GetPrimaryKey(), id)
			}()
			continue
		}
		u := userMap[id]
		chatUserRes := &models.UserJson{
			ID:       u.GetPrimaryKey(),
			Username: u.GetUsername(),
			Messages: make([]*models.MessageJson, 0),
			Unread:   0,
		}
		chatUserRes.LastChatTime = chat.AdminService.GetLastChatTime(admin.GetPrimaryKey(), u.GetPrimaryKey())
		chatUserRes.Disabled = disabled
		if _, ok := websocket.UserHub.GetConn(u.GetPrimaryKey()); ok {
			chatUserRes.Online = true
		}
		resp = append(resp, chatUserRes)
	}
	messages := messageRepo.Get([]Where{
		{
			Filed: "received_at > ?",
			Value: time.Now().Unix() - 3 * 24 * 60 * 60,
		},
		{
			Filed: "admin_id = ?",
			Value: admin.GetPrimaryKey(),
		},
		{
			Filed: "source in ?",
			Value: []int{models.SourceAdmin, models.SourceUser},
		},
	}, -1, []string{"User","Admin"}, "id desc")
	for _, u := range resp {
		for _, m := range messages {
			if m.UserId == u.ID {
				rm := m.ToJson()
				if !m.IsRead && m.Source == models.SourceUser {
					u.Unread += 1
				}
				u.Messages = append(u.Messages, rm)
			}
		}
		if _, ok := websocket.UserHub.GetConn(u.ID); ok {
			u.Online = true
		}
	}
	sort.Slice(resp, func(i, j int) bool {
		return resp[i].LastChatTime > resp[j].LastChatTime
	})
	util.RespSuccess(c, resp)
}

// 接入用户
// 分两种情况
// 一种是普通的接入
// 一种是转接的接入，转接的接入要判断转接的对象是否当前admin
func (handle *ChatHandler) AcceptUser(c *gin.Context) {
	form := &struct {
		Sid int64 `json:"sid"`
	}{}
	err := c.Bind(form)
	if err != nil {
		util.RespValidateFail(c, "invalid params")
		return
	}
	session := chatSessionRepo.First([]Where{
		{
			Filed: "id = ?",
			Value: form.Sid,
		},
	})
	if session == nil || session.CanceledAt > 0 || session.AdminId > 0 {
		util.RespNotFound(c)
		return
	}
	user := userRepo.First([]Where{
		{
			Filed: "id = ?",
			Value: session.UserId,
		},
	})
	if user == nil {
		util.RespNotFound(c)
		return
	}
	admin := auth.GetAdmin(c)
	if !admin.AccessTo(user) {
		util.RespNotFound(c)
		return
	}
	if chat.UserService.GetValidAdmin(user.GetPrimaryKey()) != 0 {
		util.RespFail(c, "user had been accepted", 10001)
		return
	}
	if session.Type == models.ChatSessionTypeTransfer {
		transferAdminId := chat.TransferService.GetUserTransferId(user.GetPrimaryKey())
		if transferAdminId == 0 {
			util.RespValidateFail(c, "transfer error ")
			return
		}
		if transferAdminId != admin.GetPrimaryKey() {
			util.RespValidateFail(c, "transfer error ")
			return
		}
		transfer := transferRepo.First([]Where{
			{
				Filed: "to_admin_id = ?",
				Value: admin.GetPrimaryKey(),
			},
			{
				Filed: "user_id = ?",
				Value: user.GetPrimaryKey(),
			},
			{
				Filed: "is_accepted = ?",
				Value: 0,
			},
		})
		if transfer.Id == 0 {
			util.RespValidateFail(c, "transfer error ")
			return
		}
		now := time.Now()
		transfer.AcceptedAt = &now
		transfer.IsAccepted = true
		transferRepo.Save(transfer)
		_ = chat.TransferService.RemoveUser(user.GetPrimaryKey())
		websocket.AdminHub.BroadcastUserTransfer(admin.GetPrimaryKey())
	}
	unSendMsg := messageRepo.GetUnSend([]Where{
		{
			Filed: "user_id = ?",
			Value: user.GetPrimaryKey(),
		},
		{
			Filed: "session_id = ?",
			Value: session.Id,
		},
	})
	sessionDuration := chat.SettingService.GetServiceSessionSecond()
	session.AcceptedAt = time.Now().Unix()
	session.AdminId = admin.GetPrimaryKey()
	chatSessionRepo.Save(session)
	_ = chat.AdminService.AddUser(admin.GetPrimaryKey(),user.GetPrimaryKey(), sessionDuration)
	now := time.Now().Unix()
	// 更新未发送的消息
	messageRepo.Update([]*repositories.Where{
		{
			Filed: "user_id = ?",
			Value: user.GetPrimaryKey(),
		},
		{
			Filed: "source = ?",
			Value: models.SourceUser,
		},
		{
			Filed: "session_id = ?",
			Value: session.Id,
		},
	}, map[string]interface{}{
		"admin_id": admin.ID,
		"send_at":  now,
	})
	messages := messageRepo.Get([]*repositories.Where{
		{
			Filed: "user_id = ?",
			Value: user.GetPrimaryKey(),
		},
		{
			Filed: "admin_id = ?",
			Value: admin.GetPrimaryKey(),
		},
		{
			Filed: "source in ?",
			Value: []int{models.SourceAdmin, models.SourceUser},
		},
	}, 20, []string{"User", "Admin"}, "id desc")
	messageLength := len(messages)
	chatUser := &models.UserJson{
		ID:           user.GetPrimaryKey(),
		Username:     user.GetUsername(),
		LastChatTime: 0,
		Messages:     make([]*models.MessageJson, messageLength, messageLength),
		Avatar: user.GetAvatarUrl(),
	}
	chatUser.Unread = len(unSendMsg)
	userConn, exist := websocket.UserHub.GetConn(user.GetPrimaryKey())
	chatUser.Online = exist
	chatUser.LastChatTime = time.Now().Unix()
	noticeMessage := &models.Message{
		UserId:     user.GetPrimaryKey(),
		AdminId:    admin.GetPrimaryKey(),
		Type:       models.TypeNotice,
		Content:    admin.GetChatName() + "为您服务",
		ReceivedAT: time.Now().Unix(),
		Source:     models.SourceSystem,
		SessionId:  session.Id,
		ReqId:      util.CreateReqId(),
	}
	if exist {
		userConn.Deliver(websocket.NewReceiveAction(noticeMessage))
	} else {
		messageRepo.Save(noticeMessage)
	}
	for index, m := range messages {
		rm := m.ToJson()
		chatUser.Messages[index] = rm
	}
	go websocket.AdminHub.BroadcastWaitingUser()
	util.RespSuccess(c, chatUser)
}

// 移除用户
func (handle *ChatHandler) RemoveUser(c *gin.Context) {
	uidStr := c.Param("id")
	admin := auth.GetAdmin(c)
	session := chatSessionRepo.First([]Where{
		{
			Filed: "user_id = ?",
			Value: uidStr,
		},
		{
			Filed: "admin_id = ?",
			Value: admin.GetPrimaryKey(),
		},
	}, "id desc")
	if session != nil {
		if session.BrokeAt == 0 {
			noticeMessage := &models.Message{
				UserId:     session.UserId,
				AdminId:    admin.GetPrimaryKey(),
				Type:       models.TypeNotice,
				Content:    admin.GetChatName() + "已断开服务",
				ReceivedAT: time.Now().Unix(),
				Source:     models.SourceSystem,
				SessionId:  session.Id,
				ReqId:      util.CreateReqId(),
			}
			userConn, exist := websocket.UserHub.GetConn(session.UserId)
			if exist {
				userConn.Deliver(websocket.NewReceiveAction(noticeMessage))
			} else {
				messageRepo.Save(noticeMessage)
			}
		}
		chat.SessionService.Close(session, true, false)
	}
	util.RespSuccess(c, nil)
}
// 已读
func (handle *ChatHandler) ReadAll(c *gin.Context) {
	form := &struct {
		Id int64
	}{}
	err := c.Bind(form)

	if err == nil {
		admin := auth.GetAdmin(c)
		wheres := []Where{
			{
				Filed: "admin_id = ?",
				Value: admin.GetPrimaryKey(),
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
		messageRepo.Update(wheres, map[string]interface{}{
			"is_read": 1,
		})
		util.RespSuccess(c, gin.H{})
	} else {
		util.RespValidateFail(c, "invalid params")
	}
}
// 获取用户信息
func (handle *ChatHandler) GetUserInfo(c *gin.Context)  {
	uidStr := c.Param("id")
	uid, err := strconv.ParseInt(uidStr, 10, 64)
	if err != nil {
		util.RespValidateFail(c, err.Error())
		return
	}
	user := userRepo.First([]Where{
		{
			Filed: "id = ?",
			Value: uid,
		},
	})
	if user == nil {
		util.RespNotFound(c)
		return
	}
	admin := auth.GetAdmin(c)
	if !admin.AccessTo(user) {
		util.RespNotFound(c)
		return
	}
	util.RespSuccess(c, gin.H{
		"username": user.GetUsername(),
		// other info
	})

}
// 转接历史消息
func (handle *ChatHandler) TransferMessages(c *gin.Context) {
	admin := auth.GetAdmin(c)
	transfer := transferRepo.First([]Where{
		{
			Filed: "to_admin_id = ?",
			Value: admin.GetPrimaryKey(),
		},
		{
			Filed: "id = ?",
			Value: c.Param("id"),
		},
	})
	if transfer == nil {
		util.RespNotFound(c)
		return
	}
	messages := messageRepo.Get([]Where{
		{
			Filed: "session_id = ?",
			Value: transfer.SessionId,
		},
	}, -1 , []string{"Admin", "User"}, "id desc")
	res := make([]*models.MessageJson, 0, len(messages))
	for _, m := range messages {
		res = append(res, m.ToJson())
	}
	util.RespSuccess(c, res)
}
// 取消转接
func (handle *ChatHandler) ChatCancelTransfer(c *gin.Context) {
	id := c.Param("id")
	admin := auth.GetAdmin(c)
	transfer := transferRepo.First([]*repositories.Where{
		{
			Filed: "to_admin_id = ?",
			Value:  admin.GetPrimaryKey(),
		},
		{
			Filed: "id = ?",
			Value: id,
		},
	})
	if transfer == nil {
		util.RespNotFound(c)
		return
	}
	if transfer.IsCanceled {
		util.RespValidateFail(c, "transfer is canceled")
		return
	}
	if transfer.IsAccepted {
		util.RespValidateFail(c, "transfer is accepted")
		return
	}
	_ = chat.TransferService.Cancel(transfer)
	websocket.AdminHub.BroadcastUserTransfer(admin.GetPrimaryKey())
	util.RespSuccess(c , gin.H{})
}
// 转接
func (handle *ChatHandler) Transfer(c *gin.Context) {
	form := &struct {
		UserId int64 `json:"user_id" binding:"required"`
		ToId int64 `json:"to_id" binding:"required,max=255"`
		Remark string `json:"remark"`
	}{}
	err := c.ShouldBind(form)
	if err != nil {
		util.RespValidateFail(c , err.Error())
		return
	}
	user := userRepo.First([]Where{
		{
			Filed: "id = ?",
			Value: form.UserId,
		},
	})
	if user == nil {
		util.RespNotFound(c)
		return
	}
	admin := auth.GetAdmin(c)
	if !admin.AccessTo(user) {
		util.RespNotFound(c)
		return
	}
	toAdmin := adminRepo.First([]Where{
		{
			Filed: "id = ?",
			Value: form.ToId,
		},
	})
	if toAdmin.ID == 0 {
		util.RespValidateFail(c , "admin_not_exist")
		return
	}
	err = chat.TransferService.Create(admin.GetPrimaryKey(), form.ToId, form.UserId, form.Remark)
	if err != nil {
		util.RespValidateFail(c , err.Error())
		return
	}
	go websocket.AdminHub.BroadcastUserTransfer(form.ToId)
	util.RespSuccess(c, gin.H{})
}
// 聊天图片
func (handle *ChatHandler) Image(c *gin.Context) {
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
