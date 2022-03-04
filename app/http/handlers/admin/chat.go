package admin

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
	"ws/app/chat"
	"ws/app/contract"
	"ws/app/http/requests"
	"ws/app/models"
	"ws/app/repositories"
	"ws/app/resource"
	"ws/app/util"
	"ws/app/websocket"
)

type ChatHandler struct {
}

// GetHistorySession 查看历史对话
func (handle *ChatHandler) GetHistorySession(c *gin.Context) {
	useId := c.Param("uid")
	sessions := chatSessionRepo.Get([]Where{
		{
			Filed: "user_id = ?",
			Value: useId,
		},
		{
			Filed: "group_id = ?",
			Value: requests.GetAdmin(c).GetGroupId(),
		},
		{
			Filed: "admin_id > ?",
			Value: 0,
		},
	}, -1, []string{"Admin","User"}, []string{"id desc"})
	resp := make([]*resource.ChatSession, len(sessions))
	for i, session := range sessions {
		resp[i] = session.ToJson()
	}
	util.RespSuccess(c, resp)
}

// GetHistoryMessage 获取消息
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
	admin := requests.GetAdmin(c)
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
	messages := messageRepo.Get(wheres, 20, []string{"User", "Admin"}, []string{"id desc"})
	res := make([]*resource.Message, 0)
	for _, m := range messages {
		res = append(res, m.ToJson())
	}
	util.RespSuccess(c, res)
}

// GetReqId 获取reqId
func (handle *ChatHandler) GetReqId(c *gin.Context)  {
	admin := requests.GetAdmin(c)
	util.RespSuccess(c, gin.H{
		"reqId" : admin.GetReqId(),
	})
}

// ChatUserList 聊天用户列表
func (handle *ChatHandler) ChatUserList(c *gin.Context) {
	admin := requests.GetAdmin(c)
	ids, times := chat.AdminService.GetUsersWithLimitTime(admin.GetPrimaryKey())
	users := userRepo.Get([]Where{
		{
			Filed: "id in ?",
			Value: ids,
		},
	}, -1, []string{}, []string{})
	resp := make([]*resource.User, 0, len(users))
	userMap := make(map[int64]contract.User)
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
		chatUserRes := &resource.User{
			ID:       u.GetPrimaryKey(),
			Username: u.GetUsername(),
			Messages: make([]*resource.Message, 0),
			Unread:   0,
		}
		chatUserRes.LastChatTime = chat.AdminService.GetLastChatTime(admin.GetPrimaryKey(), u.GetPrimaryKey())
		chatUserRes.Disabled = disabled
		if _, ok := websocket.UserManager.GetConn(u); ok {
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
	}, -1, []string{"User","Admin"}, []string{"id desc"})
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
	}
	util.RespSuccess(c, resp)
}

// AcceptUser 接入用户
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
	}, []string{})
	if session == nil || session.CanceledAt > 0 || session.AdminId > 0 {
		util.RespNotFound(c)
		return
	}
	user := userRepo.First([]Where{
		{
			Filed: "id = ?",
			Value: session.UserId,
		},
	}, []string{})
	if user == nil {
		util.RespNotFound(c)
		return
	}
	u := requests.GetAdmin(c)
	admin, _ := u.(*models.Admin)
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
		}, []string{})
		if transfer == nil {
			util.RespValidateFail(c, "transfer error ")
			return
		}
		now := time.Now()
		transfer.AcceptedAt = now.Unix()
		transfer.IsAccepted = true
		transferRepo.Save(transfer)
		_ = chat.TransferService.RemoveUser(user.GetPrimaryKey())
		websocket.AdminManager.PublishTransfer(admin)
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
	session.AcceptedAt = time.Now().Unix()
	session.AdminId = admin.GetPrimaryKey()
	chatSessionRepo.Save(session)
	_ = chat.AdminService.AddUser(admin, user)
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
		"admin_id": admin.GetPrimaryKey(),
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
	}, 20, []string{"User", "Admin"}, []string{"id desc"})
	messageLength := len(messages)
	chatUser := &resource.User{
		ID:           user.GetPrimaryKey(),
		Username:     user.GetUsername(),
		LastChatTime: 0,
		Messages:     make([]*resource.Message, messageLength, messageLength),
		Avatar: user.GetAvatarUrl(),
	}
	chatUser.Unread = len(unSendMsg)
	chatUser.LastChatTime = time.Now().Unix()
	noticeMessage := messageRepo.NewNotice(session, admin.GetChatName() + "为您服务")
	messageRepo.Save(noticeMessage)
	websocket.UserManager.DeliveryMessage(noticeMessage, false)
	for index, m := range messages {
		rm := m.ToJson()
		chatUser.Messages[index] = rm
	}
	go websocket.AdminManager.PublishWaitingUser(user.GetGroupId())
	go websocket.UserManager.PublishWaitingCount(user.GetGroupId())
	util.RespSuccess(c, chatUser)
}

// RemoveUser 移除用户
func (handle *ChatHandler) RemoveUser(c *gin.Context) {
	uidStr := c.Param("id")
	u := requests.GetAdmin(c)
	admin, _ := u.(*models.Admin)
	session := chatSessionRepo.First([]Where{
		{
			Filed: "user_id = ?",
			Value: uidStr,
		},
		{
			Filed: "admin_id = ?",
			Value: admin.GetPrimaryKey(),
		},
	}, []string{"id desc"})
	if session != nil {
		if session.BrokeAt == 0 {
			noticeMessage := messageRepo.NewNotice(session, admin.GetChatName() + "已断开服务")
			messageRepo.Save(noticeMessage)
			websocket.UserManager.DeliveryMessage(noticeMessage, false)
		}
		chat.SessionService.Close(session.Id, true, false)
	}
	util.RespSuccess(c, nil)
}

// ReadAll 已读
func (handle *ChatHandler) ReadAll(c *gin.Context) {
	form := &struct {
		Id int64 `json:"id"`
		MsgId int64 `json:"msg_id" binding:"-"`
	}{}
	err := c.Bind(form)
	admin := requests.GetAdmin(c)
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
	if err == nil {
		if form.MsgId > 0 {
			wheres = append(wheres, &repositories.Where{
				Filed: "id <= ?",
				Value: form.MsgId,
			})
		}
		messageRepo.Update(wheres, map[string]interface{}{
			"is_read": 1,
		})
		util.RespSuccess(c, gin.H{})
	} else {
		util.RespValidateFail(c, "invalid params")
	}
}

// GetUserInfo 获取用户信息
func (handle *ChatHandler) GetUserInfo(c *gin.Context)  {
	uidStr := c.Param("id")
	uid, err := strconv.ParseInt(uidStr, 10, 64)
	if err != nil {
		util.RespValidateFail(c, err.Error())
		return
	}
	admin := requests.GetAdmin(c)
	user := userRepo.First([]Where{
		{
			Filed: "id = ?",
			Value: uid,
		},
		{
			Filed: "group_id = ?",
			Value: admin.GetGroupId(),
		},
	}, []string{})
	if user == nil {
		util.RespNotFound(c)
		return
	}
	if !admin.AccessTo(user) {
		util.RespNotFound(c)
		return
	}
	util.RespSuccess(c, gin.H{
		"username": user.GetUsername(),
		// other info
	})

}

// TransferMessages 转接历史消息
func (handle *ChatHandler) TransferMessages(c *gin.Context) {
	admin := requests.GetAdmin(c)
	transfer := transferRepo.First([]Where{
		{
			Filed: "to_admin_id = ?",
			Value: admin.GetPrimaryKey(),
		},
		{
			Filed: "id = ?",
			Value: c.Param("id"),
		},
	}, []string{})
	if transfer == nil {
		util.RespNotFound(c)
		return
	}
	messages := messageRepo.Get([]Where{
		{
			Filed: "session_id = ?",
			Value: transfer.SessionId,
		},
	}, -1 , []string{"Admin", "User"}, []string{"id desc"})
	res := make([]*resource.Message, 0, len(messages))
	for _, m := range messages {
		res = append(res, m.ToJson())
	}
	util.RespSuccess(c, res)
}

// CancelTransfer 取消转接
func (handle *ChatHandler) CancelTransfer(c *gin.Context) {
	id := c.Param("id")
	admin := requests.GetAdmin(c)
	transfer := transferRepo.First([]*repositories.Where{
		{
			Filed: "to_admin_id = ?",
			Value:  admin.GetPrimaryKey(),
		},
		{
			Filed: "id = ?",
			Value: id,
		},
	}, []string{})
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
	websocket.AdminManager.PublishTransfer(admin)
	util.RespSuccess(c , gin.H{})
}

// Transfer 转接
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
	}, []string{})
	if user == nil {
		util.RespNotFound(c)
		return
	}
	admin := requests.GetAdmin(c)
	if !admin.AccessTo(user) {
		util.RespNotFound(c)
		return
	}
	toAdmin := adminRepo.First([]Where{
		{
			Filed: "id = ?",
			Value: form.ToId,
		},
	}, []string{})
	if toAdmin.ID == 0 {
		util.RespValidateFail(c , "admin_not_exist")
		return
	}
	err = chat.TransferService.Create(admin.GetPrimaryKey(), form.ToId, form.UserId, form.Remark)
	if err != nil {
		util.RespValidateFail(c , err.Error())
		return
	}
	go websocket.AdminManager.PublishTransfer(toAdmin)
	util.RespSuccess(c, gin.H{})
}

