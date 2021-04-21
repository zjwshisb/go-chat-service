package http

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
	"ws/db"
	"ws/hub"
	"ws/models"
	"ws/util"
)

// 接入用户
func AcceptUser(c *gin.Context) {
	form := &struct {
		Uid int64
	}{}
	err := c.Bind(form)
	if err == nil {
		ui, _ := c.Get("user")
		serverUser := ui.(*models.ServerUser)
		sClient, exist := hub.Hub.Server.GetClient(serverUser.ID)
		if !exist {
			util.RespFail(c,"请先登录", 500)
			return
		}
		user, err := sClient.Accept(form.Uid)
		if err != nil {
			util.RespFail(c, err.Error(), 500)
			return
		}
		unSendMsg := user.GetUnSendMsg()
		now := time.Now().Unix()
		// 更新未发送的消息
		db.Db.Table("messages").
			Where("user_id = ?", form.Uid).
			Where("service_id = ?", 0).Updates(map[string]interface{}{
				"service_id": serverUser.ID,
				"send_at":now,
		})
		messages := make([]models.Message, 0)
		db.Db.Where("user_id = ?", form.Uid).
			Where("service_id = ?", serverUser.ID).
			Where("received_at >= ?", now - 2 * 24 * 60 * 60).Find(&messages)
		for _, m:= range messages {
			m.IsSuccess = true
		}
		chatUser := models.ChatUser{
			ID: user.ID,
			Username: user.Username,
			Online: true,
			Disabled: false,
			Messages: messages,
			Unread: len(unSendMsg),
			LastChatTime: time.Now().Unix(),
		}
		util.RespSuccess(c, chatUser)
	} else {
		util.RespValidateFail(c , "验证不通过")
	}
}
// 移除用户
func RemoveUser(c *gin.Context) {
	uidStr := c.Param("id")
	uid, err := strconv.ParseInt(uidStr, 10, 64)
	if err == nil {
		ui, _ := c.Get("user")
		user := ui.(*models.ServerUser)
		_ = user.RemoveChatUser(uid)
		if uClient, exist := hub.Hub.User.AcceptedClient.GetClient(uid); exist {
			uClient.RemoveServerId()
			if uClient.ServerId == user.ID {
				hub.Hub.User.Change2waiting(uClient)
			}
		}
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
		server := ui.(*models.ServerUser)
		db.Db.Model(&models.Message{}).
			Where("service_id = ?" , server.ID).
			Where("user_id = ?", form.Id).
			Update("is_read", 1)
		util.RespSuccess(c, gin.H{})
	} else {
		util.RespValidateFail(c, "invalid params")
	}
}
