package admin

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
	"ws/app/auth"
	"ws/app/chat"
	"ws/app/databases"
	"ws/app/file"
	"ws/app/models"
	"ws/app/repositories"
	"ws/app/util"
	"ws/app/websocket"
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
	admin := auth.GetAdmin(c)
	chatIds, _ := chat.GetAdminUserIds(admin.GetPrimaryKey())
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
	wheres := []*repositories.Where{
		{
			Filed: "admin_id = ?",
			Value: admin.GetPrimaryKey(),
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
	messages := repositories.GetMessages(wheres, 20, []string{"User", "Admin"})
	res := make([]*models.MessageJson, 0)
	for _, m := range messages {
		res = append(res, m.ToJson())
	}
	util.RespSuccess(c, res)
}

// 聊天用户列表
func ChatUserList(c *gin.Context) {
	admin := auth.GetAdmin(c)
	ids, times := chat.GetAdminUserIds(admin.GetPrimaryKey())
	users := repositories.GetUserByIds(ids)
	resp := make([]*models.UserJson, 0, len(users))
	userMap := make(map[int64]auth.User)
	for _, user := range users {
		userMap[user.GetPrimaryKey()] = user
	}
	for index, id := range ids {
		u := userMap[id]
		chatUserRes := &models.UserJson{
			ID:       u.GetPrimaryKey(),
			Username: u.GetUsername(),
			Messages: make([]*models.MessageJson, 0),
			Unread:   0,
		}
		limitTime := times[index]
		chatUserRes.LastChatTime = chat.GetAdminUserLastChatTime(u.GetPrimaryKey(), admin.GetPrimaryKey() )
		chatUserRes.Disabled = limitTime <= time.Now().Unix()
		if _, ok := websocket.UserHub.GetConn(u.GetPrimaryKey()); ok {
			chatUserRes.Online = true

		}
		resp = append(resp, chatUserRes)
	}
	messages := repositories.GetMessages([]*repositories.Where{
		{
			Filed: "received_at > ?",
			Value: time.Now().Unix() - 3 * 24 * 60 * 60,
		},
		{
			Filed: "admin_id = ?",
			Value: admin.GetPrimaryKey(),
		},
	}, -1, []string{"User","Admin"})
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
	util.RespSuccess(c, resp)
}

// 接入用户
// 分两种情况
// 一种是普通的接入
// 一种是转接的接入，转接的接入要判断转接的对象是否当前admin
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
	if chat.GetUserLastAdminId(user.GetPrimaryKey()) != 0 {
		util.RespFail(c, "user had been accepted", 10001)
		return
	}
	session := chat.GetSession(user.GetPrimaryKey(), 0)
	if session == nil {
		util.RespValidateFail(c , "chat session error")
		return
	}
	admin := auth.GetAdmin(c)
	if session.Type == models.ChatSessionTypeTransfer {
		transferAdminId := chat.GetUserTransferId(user.GetPrimaryKey())
		if transferAdminId == 0 {
			util.RespValidateFail(c, "transfer error ")
			return
		}
		if transferAdminId != admin.GetPrimaryKey() {
			util.RespValidateFail(c, "transfer error ")
			return
		}
		transfer := &models.ChatTransfer{}
		databases.Db.Where("to_admin_id = ?", admin.GetPrimaryKey()).
			Where("user_id = ?", user.GetPrimaryKey()).
			Where("is_accepted = ?" ,0).Find(transfer)
		if transfer.Id == 0 {
			util.RespValidateFail(c, "transfer error ")
			return
		}
		now := time.Now()
		transfer.AcceptedAt = &now
		transfer.IsAccepted = true
		databases.Db.Save(transfer)
		_ = chat.RemoveTransfer(user.GetPrimaryKey())
		websocket.AdminHub.BroadcastUserTransfer(admin.GetPrimaryKey())
	}
	unSendMsg := repositories.GetUnSendMessage(
		&repositories.Where{
			Filed: "user_id = ?",
			Value: user.GetPrimaryKey(),
		},
		&repositories.Where{
			Filed: "source = ?",
			Value: models.SourceUser,
		},
	)
	sessionDuration := chat.GetServiceSessionSecond()
	session.AcceptedAt = time.Now().Unix()
	session.AdminId = admin.GetPrimaryKey()
	session.BrokeAt = time.Now().Unix() + sessionDuration
	databases.Db.Save(session)
	_ = chat.SetUserAdminId(user.GetPrimaryKey(), admin.GetPrimaryKey(), sessionDuration)
	now := time.Now().Unix()
	// 更新未发送的消息
	repositories.UpdateMessages([]*repositories.Where{
		{
			Filed: "user_id = ?",
			Value: user.GetPrimaryKey(),
		},
		{
			Filed: "source = ?",
			Value: models.SourceUser,
		},
		{
			Filed: "admin_id = ?",
			Value: 0,
		},
	}, map[string]interface{}{
		"admin_id": admin.ID,
		"send_at":    now,
		"session_id": session.Id,
	})
	messages := repositories.GetMessages([]*repositories.Where{
		{
			Filed: "user_id = ?",
			Value: user.GetPrimaryKey(),
		},
		{
			Filed: "admin_id = ?",
			Value: admin.GetPrimaryKey(),
		},
	}, 20, []string{"User", "Admin"})
	messageLength := len(messages)
	chatUser := &models.UserJson{
		ID:           user.GetPrimaryKey(),
		Username:     user.GetUsername(),
		LastChatTime: 0,
		Messages:     make([]*models.MessageJson, messageLength, messageLength),
	}
	chatUser.Unread = len(unSendMsg)
	_, exist = websocket.UserHub.GetConn(user.GetPrimaryKey())
	chatUser.Online = exist
	chatUser.LastChatTime = time.Now().Unix()
	for index, m := range messages {
		rm := m.ToJson()
		chatUser.Messages[index] = rm
	}
	go websocket.AdminHub.BroadcastWaitingUser()
	util.RespSuccess(c, chatUser)
}

// 移除用户
func RemoveUser(c *gin.Context) {
	uidStr := c.Param("id")
	uid, err := strconv.ParseInt(uidStr, 10, 64)
	admin := auth.GetAdmin(c)
	if err == nil {
		_ = chat.RemoveUserAdminId(uid, admin.GetPrimaryKey())
	}
	record := &models.ChatSession{}
	databases.Db.Where("user_id = ?", uidStr).
		Where("admin_id = ?" , admin.GetPrimaryKey()).
		Order("id desc").First(record)
	if record.Id > 0 {
		record.BrokeAt = time.Now().Unix()
		databases.Db.Save(record)
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
		admin := auth.GetAdmin(c)
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
		repositories.UpdateMessages(wheres, map[string]interface{}{
			"is_read": 1,
		})
		util.RespSuccess(c, gin.H{})
	} else {
		util.RespValidateFail(c, "invalid params")
	}
}
func GetUserInfo(c *gin.Context)  {
	uidStr := c.Param("id")
	uid, err := strconv.ParseInt(uidStr, 10, 64)
	if err != nil {
		util.RespValidateFail(c, err.Error())
		return
	}
	user, exist := repositories.GetUserById(uid)
	if !exist {
		util.RespNotFound(c)
		return
	}
	util.RespSuccess(c, gin.H{
		"username": user.GetUsername(),
		// other info
	})

}
func TransferMessages(c *gin.Context) {
	transfer := &models.ChatTransfer{}
	admin := auth.GetAdmin(c)
	databases.Db.Where("to_admin_id = ?", admin.GetPrimaryKey()).Find(transfer, c.Param("id"))
	if transfer.Id == 0 {
		util.RespNotFound(c)
		return
	}
	messages := make([]*models.Message, 0)
	databases.Db.Where("session_id = ?", transfer.SessionId).Preload("Admin").
		Preload("User").Find(&messages)
	res := make([]*models.MessageJson, 0, len(messages))
	for _, m := range messages {
		res = append(res, m.ToJson())
	}
	util.RespSuccess(c, res)
}
// 转发
func Transfer(c *gin.Context) {
	admin := auth.GetAdmin(c)
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
	toAdmin := &models.Admin{}
	databases.Db.Find(toAdmin, form.ToId)
	if toAdmin.ID == 0 {
		util.RespValidateFail(c , "admin_not_exist")
		return
	}
	err = chat.Transfer(admin.GetPrimaryKey(), form.ToId, form.UserId, form.Remark)
	if err != nil {
		util.RespValidateFail(c , err.Error())
		return
	}
	go websocket.AdminHub.BroadcastUserTransfer(form.ToId)
	util.RespSuccess(c, gin.H{})
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
