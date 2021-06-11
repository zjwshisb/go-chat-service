package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
	"ws/configs"
	"ws/internal/databases"
	"ws/internal/file"
	"ws/internal/models"
	"ws/internal/resources"
	"ws/internal/websocket"
	"ws/util"
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
	user := getUser(c)
	var messages []*models.Message
	query := databases.Db.Where("service_id = ?", user.ID).
		Where("user_id = ?", uid)
	midStr, exist := c.GetQuery("mid")
	if exist {
		mid, err = strconv.ParseInt(midStr, 10, 64)
		if err == nil {
			query.Where("id < ?", mid)
		}
	}
	query.Order("id desc").Limit(20).Find(&messages)
	res := make([]*resources.Message, 0)
	for _, m := range messages {
		res = append(res, resources.NewMessage(*m))
	}
	util.RespSuccess(c, res)
}

// 聊天用户列表
func ChatUserList(c *gin.Context) {
	user := getUser(c)
	chatUsers := user.GetChatUsers()
	cmd := databases.Redis.ZRevRangeWithScores(context.Background(), user.ChatUsersKey(), 0, -1)
	resp := make([]*resources.ChatUser, 0, len(cmd.Val()))
	if cmd.Err() == redis.Nil {
		util.RespSuccess(c, resp)
	} else {
		for _, z := range cmd.Val() {
			id, err := strconv.ParseInt(z.Member.(string), 10, 64)
			if err == nil {
				for _, user := range chatUsers {
					if user.ID == id {
						chatUserRes := &resources.ChatUser{
							ID:           user.ID,
							Username:     user.Username,
							Messages:     make([]*resources.Message, 0),
							Unread:       0,
						}
						chatUserRes.LastChatTime = int64(z.Score)
						resp = append(resp, chatUserRes)
					}
				}
			}
		}
	}

	last := time.Now().Unix() - configs.App.ChatSessionDuration * 24 * 60 * 60
	var messages []*models.Message
	databases.Db.Preload("ServerUser").
		Preload("User").
		Where("received_at > ?", last).
		Where("service_id = ?", user.ID).
		Find(&messages)
	for _, u := range resp {
		for _, m := range messages {
			if m.UserId == u.ID {
				rm := resources.NewMessage(*m)
				rm.IsSuccess = true
				if !m.IsRead && !m.IsServer {
					u.Unread += 1
				}
				u.Messages = append(u.Messages, rm)
			}
		}
		u.Disabled = !user.CheckChatUserLegal(u.ID)
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
	ui, _ := c.Get("user")
	serverUser := ui.(*models.ServiceUser)
	var user models.User
	databases.Db.Where("id = ?", form.Uid).First(&user)
	if user.ID == 0 {
		util.RespValidateFail(c, "invalid params")
		return
	}
	if user.GetLastServiceId() != 0 {
		util.RespFail(c, "use had been accepted", 10001)
		return
	}
	unSendMsg := user.GetUnSendMsg()
	now := time.Now().Unix()
	// 更新未发送的消息
	databases.Db.Table("messages").
		Where("user_id = ?", form.Uid).
		Where("service_id = ?", 0).Updates(map[string]interface{}{
		"service_id": serverUser.ID,
		"send_at":    now,
	})
	messages := make([]*models.Message, 0)
	databases.Db.Where("user_id = ?", form.Uid).
		Where("service_id = ?", serverUser.ID).
		Limit(20).
		Find(&messages)

	_ = serverUser.UpdateChatUser(user.ID)
	_ = user.SetServiceId(serverUser.ID)

	chatUser := &resources.ChatUser{
		ID:           user.ID,
		Username:     user.Username,
		LastChatTime: 0,
		Messages:     make([]*resources.Message, 0, len(messages)),
	}
	chatUser.Unread = len(unSendMsg)
	_, exist := websocket.UserHub.GetConn(user.ID)
	chatUser.Online = exist
	chatUser.LastChatTime = time.Now().Unix()
	for _, m := range messages {
		rm := resources.NewMessage(*m)
		chatUser.Messages = append(chatUser.Messages, rm)
	}
	go websocket.ServiceHub.BroadcastWaitingUser()
	util.RespSuccess(c, chatUser)
}

// 移除用户
func RemoveUser(c *gin.Context) {
	uidStr := c.Param("id")
	uid, err := strconv.ParseInt(uidStr, 10, 64)
	if err == nil {
		ui, _ := c.Get("user")
		serviceUser := ui.(*models.ServiceUser)
		_ = serviceUser.RemoveChatUser(uid)
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
		ui, _ := c.Get("user")
		server := ui.(*models.ServiceUser)
		databases.Db.Model(&models.Message{}).
			Where("service_id = ?", server.ID).
			Where("user_id = ?", form.Id).
			Update("is_read", 1)
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
