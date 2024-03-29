package admin

import (
	"github.com/duke-git/lancet/v2/random"
	"github.com/duke-git/lancet/v2/slice"
	"strconv"
	"time"
	"ws/app/chat"
	"ws/app/contract"
	"ws/app/http/requests"
	"ws/app/http/responses"
	"ws/app/http/websocket"
	"ws/app/models"
	"ws/app/repositories"
	"ws/app/resource"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
}

// GetHistorySession 查看历史对话
func (handle *ChatHandler) GetHistorySession(c *gin.Context) {
	useId := c.Param("uid")
	sessions := repositories.ChatSessionRepo.Get([]*repositories.Where{
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
	}, -1, []string{"Admin", "User"}, []string{"id desc"})
	resp := make([]*resource.ChatSession, len(sessions))
	for i, session := range sessions {
		resp[i] = session.ToJson()
	}

	responses.RespSuccess(c, resp)
}

// GetHistoryMessage 获取消息
func (handle *ChatHandler) GetHistoryMessage(c *gin.Context) {
	var uid int64
	var mid int64
	var err error
	uidStr, exist := c.GetQuery("uid")
	if !exist {
		responses.RespValidateFail(c, "invalid params")
		return
	}
	uid, err = strconv.ParseInt(uidStr, 10, 64)
	if err != nil {
		responses.RespValidateFail(c, "invalid params")
		return
	}
	admin := requests.GetAdmin(c)
	chatIds, _ := chat.AdminService.GetUsersWithLimitTime(admin.GetPrimaryKey())
	userExist := false
	for _, chatId := range chatIds {
		if chatId == uid {
			userExist = true
		}
	}
	if !userExist {
		responses.RespValidateFail(c, "invalid params")
		return
	}
	wheres := []*repositories.Where{
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
	messages := repositories.MessageRepo.Get(wheres, 20, []string{"User", "Admin"}, []string{"id desc"})
	res := make([]*resource.Message, len(messages), len(messages))
	msgIds := make([]int64, len(messages), len(messages))
	for i, m := range messages {
		res[i] = m.ToJson()
		msgIds[i] = m.Id
	}
	responses.RespSuccess(c, res)
	updateWhere := []*repositories.Where{
		{
			Filed: "id in ?",
			Value: msgIds,
		},
		{
			Filed: "send_at = ?",
			Value: 0,
		},
		{
			Filed: "source in ?",
			Value: []int{models.SourceSystem, models.SourceUser},
		},
	}
	repositories.MessageRepo.Update(updateWhere, map[string]interface{}{
		"send_at": time.Now().Unix(),
	})
}

// GetReqId 获取reqId
func (handle *ChatHandler) GetReqId(c *gin.Context) {
	responses.RespSuccess(c, gin.H{
		"reqId": random.RandString(20),
	})
}

// ChatUserList 聊天用户列表
func (handle *ChatHandler) ChatUserList(c *gin.Context) {
	admin := requests.GetAdmin(c)
	ids, times := chat.AdminService.GetUsersWithLimitTime(admin.GetPrimaryKey())
	users := repositories.UserRepo.Get([]*repositories.Where{
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
	messages := repositories.MessageRepo.Get([]*repositories.Where{
		{
			Filed: "received_at > ?",
			Value: time.Now().Unix() - 3*24*60*60,
		},
		{
			Filed: "admin_id = ?",
			Value: admin.GetPrimaryKey(),
		},
		{
			Filed: "source in ?",
			Value: []int{models.SourceAdmin, models.SourceUser},
		},
	}, -1, []string{"User", "Admin"}, []string{"id desc"})
	messageIds := make([]int64, len(messages), len(messages))
	for _, u := range resp {
		for i, m := range messages {
			messageIds[i] = m.Id
			if m.UserId == u.ID {
				rm := m.ToJson()
				if !m.IsRead && m.Source == models.SourceUser {
					u.Unread += 1
				}
				u.Messages = append(u.Messages, rm)
			}
		}
	}
	responses.RespSuccess(c, resp)
	updateWhere := []*repositories.Where{
		{
			Filed: "id in ?",
			Value: messageIds,
		},
		{
			Filed: "source in ?",
			Value: []int{models.SourceUser, models.SourceSystem},
		},
		{
			Filed: "send_at = ?",
			Value: 0,
		},
	}
	repositories.MessageRepo.Update(updateWhere, map[string]interface{}{
		"send_at": time.Now().Unix(),
	})
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
		responses.RespValidateFail(c, "invalid params")
		return
	}
	session := repositories.ChatSessionRepo.First([]*repositories.Where{
		{
			Filed: "id = ?",
			Value: form.Sid,
		},
	}, []string{})
	if session == nil || session.CanceledAt > 0 || session.AdminId > 0 {
		responses.RespNotFound(c)
		return
	}
	user := repositories.UserRepo.First([]*repositories.Where{
		{
			Filed: "id = ?",
			Value: session.UserId,
		},
	}, []string{})
	if user == nil {
		responses.RespNotFound(c)
		return
	}
	u := requests.GetAdmin(c)
	admin, _ := u.(*models.Admin)
	if !admin.AccessTo(user) {
		responses.RespNotFound(c)
		return
	}
	if chat.UserService.GetValidAdmin(user.GetPrimaryKey()) != 0 {
		responses.RespFail(c, "user had been accepted", 10001)
		return
	}
	if session.Type == models.ChatSessionTypeTransfer {
		transferAdminId := chat.TransferService.GetUserTransferId(user.GetPrimaryKey())
		if transferAdminId == 0 {
			responses.RespValidateFail(c, "transfer error ")
			return
		}
		if transferAdminId != admin.GetPrimaryKey() {
			responses.RespValidateFail(c, "transfer error ")
			return
		}
		transfer := repositories.TransferRepo.First([]*repositories.Where{
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
			responses.RespValidateFail(c, "transfer error ")
			return
		}
		now := time.Now()
		transfer.AcceptedAt = now.Unix()
		transfer.IsAccepted = true
		repositories.TransferRepo.Save(transfer)
		_ = chat.TransferService.RemoveUser(user.GetPrimaryKey())
		websocket.AdminManager.NoticeUserTransfer(admin)
	}
	unSendMsg := repositories.MessageRepo.GetUnSend([]*repositories.Where{
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
	repositories.ChatSessionRepo.Save(session)
	_ = chat.AdminService.AddUser(admin, user)
	now := time.Now().Unix()
	// 更新未发送的消息
	repositories.MessageRepo.Update([]*repositories.Where{
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
	messages := repositories.MessageRepo.Get([]*repositories.Where{
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
		Avatar:       user.GetAvatarUrl(),
	}
	chatUser.Unread = len(unSendMsg)
	chatUser.LastChatTime = time.Now().Unix()
	chatUser.Online = websocket.UserManager.IsOnline(user)
	noticeMessage := repositories.MessageRepo.NewNotice(session, admin.GetChatName()+"为您服务")

	repositories.MessageRepo.Save(noticeMessage)
	websocket.UserManager.DeliveryMessage(noticeMessage, false)
	for index, m := range messages {
		rm := m.ToJson()
		chatUser.Messages[index] = rm
	}
	go websocket.AdminManager.BroadcastWaitingUser(user.GetGroupId())
	go websocket.UserManager.BroadcastQueueLocation(user.GetGroupId())
	responses.RespSuccess(c, chatUser)
}

// RemoveUser 移除用户
func (handle *ChatHandler) RemoveUser(c *gin.Context) {
	uidStr := c.Param("id")
	u := requests.GetAdmin(c)
	admin, _ := u.(*models.Admin)
	session := repositories.ChatSessionRepo.First([]*repositories.Where{
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
			noticeMessage := repositories.MessageRepo.NewNotice(session, admin.GetChatName()+"已断开服务")
			repositories.MessageRepo.Save(noticeMessage)
			websocket.UserManager.DeliveryMessage(noticeMessage, false)
		}
		chat.SessionService.Close(session.Id, true, false)
	}
	responses.RespSuccess(c, nil)
}

// ReadAll 已读
func (handle *ChatHandler) ReadAll(c *gin.Context) {
	form := &struct {
		Id    int64 `json:"id"`
		MsgId int64 `json:"msg_id" binding:"-"`
	}{}
	err := c.Bind(form)
	admin := requests.GetAdmin(c)
	wheres := []*repositories.Where{
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
		repositories.MessageRepo.Update(wheres, map[string]interface{}{
			"is_read": 1,
		})
		responses.RespSuccess(c, gin.H{})
	} else {
		responses.RespValidateFail(c, "invalid params")
	}
}

// GetUserInfo 获取用户信息
func (handle *ChatHandler) GetUserInfo(c *gin.Context) {
	uidStr := c.Param("id")
	uid, err := strconv.ParseInt(uidStr, 10, 64)
	if err != nil {
		responses.RespValidateFail(c, err.Error())
		return
	}
	admin := requests.GetAdmin(c)
	user := repositories.UserRepo.First([]*repositories.Where{
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
		responses.RespNotFound(c)
		return
	}
	if !admin.AccessTo(user) {
		responses.RespNotFound(c)
		return
	}
	responses.RespSuccess(c, gin.H{
		"username": user.GetUsername(),
		// other info
	})

}

// TransferMessages 转接历史消息
func (handle *ChatHandler) TransferMessages(c *gin.Context) {
	admin := requests.GetAdmin(c)
	transfer := repositories.TransferRepo.First([]*repositories.Where{
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
		responses.RespNotFound(c)
		return
	}
	messages := repositories.MessageRepo.Get([]*repositories.Where{
		{
			Filed: "session_id = ?",
			Value: transfer.SessionId,
		},
	}, -1, []string{"Admin", "User"}, []string{"id desc"})
	res := slice.Map(messages, func(index int, s *models.Message) *resource.Message {
		return s.ToJson()
	})
	responses.RespSuccess(c, res)
}

// CancelTransfer 取消转接
func (handle *ChatHandler) CancelTransfer(c *gin.Context) {
	id := c.Param("id")
	admin := requests.GetAdmin(c)
	transfer := repositories.TransferRepo.First([]*repositories.Where{
		{
			Filed: "to_admin_id = ?",
			Value: admin.GetPrimaryKey(),
		},
		{
			Filed: "id = ?",
			Value: id,
		},
	}, []string{})
	if transfer == nil {
		responses.RespNotFound(c)
		return
	}
	if transfer.IsCanceled {
		responses.RespValidateFail(c, "transfer is canceled")
		return
	}
	if transfer.IsAccepted {
		responses.RespValidateFail(c, "transfer is accepted")
		return
	}
	_ = chat.TransferService.Cancel(transfer)
	websocket.AdminManager.NoticeUserTransfer(admin)
	responses.RespSuccess(c, gin.H{})
}

// Transfer 转接
func (handle *ChatHandler) Transfer(c *gin.Context) {
	form := &struct {
		UserId int64  `json:"user_id" binding:"required"`
		ToId   int64  `json:"to_id" binding:"required,max=255"`
		Remark string `json:"remark"`
	}{}
	err := c.ShouldBind(form)
	admin := requests.GetAdmin(c)
	if err != nil {
		responses.RespValidateFail(c, err.Error())
		return
	}
	user := repositories.UserRepo.First([]*repositories.Where{
		{
			Filed: "group_id =?",
			Value: admin.GetGroupId(),
		},
		{
			Filed: "id = ?",
			Value: form.UserId,
		},
	}, []string{})
	if user == nil {
		responses.RespNotFound(c)
		return
	}
	if !admin.AccessTo(user) {
		responses.RespNotFound(c)
		return
	}
	toAdmin := repositories.AdminRepo.First([]*repositories.Where{
		{
			Filed: "group_id =?",
			Value: admin.GetGroupId(),
		},
		{
			Filed: "id = ?",
			Value: form.ToId,
		},
	}, []string{})
	if toAdmin.ID == 0 {
		responses.RespValidateFail(c, "admin_not_exist")
		return
	}
	err = chat.TransferService.Create(admin.GetPrimaryKey(), form.ToId, form.UserId, form.Remark)
	if err != nil {
		responses.RespValidateFail(c, err.Error())
		return
	}
	go websocket.AdminManager.NoticeUserTransfer(toAdmin)
	responses.RespSuccess(c, gin.H{})
}
