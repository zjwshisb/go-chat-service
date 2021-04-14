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
func Accept(c *gin.Context) {
	form := &struct {
		Uid int64
	}{}
	err := c.Bind(form)
	if err == nil {
		ui, _ := c.Get("user")
		serverUser := ui.(*models.ServerUser)
		sClient, exist := hub.Hub.Server.GetClient(serverUser.ID)
		if !exist {
			c.JSON(200, util.RespFail("请先登录", 500))
			return
		}
		user, err := sClient.Accept(form.Uid)
		if err != nil {
			c.JSON(200, util.RespFail(err.Error(), 500))
			return
		}
		unreadMsg := user.GetUnSendMsg()
		messages := make([]models.Message, 0)
		for _, m := range unreadMsg {
			db.Db.Model(&m).Update("service_id", serverUser.ID)
			messages = append(messages, m)
		}
		chatUser := models.ChatUser{
			ID: user.ID,
			Username: user.Username,
			Online: true,
			Disabled: false,
			Messages: messages,
			Unread: len(unreadMsg),
			LastChatTime: time.Now().Unix(),
		}
		c.JSON(200, util.RespSuccess(chatUser))
	} else {
		c.JSON(200, util.RespFail("测试", 500))
	}
}

func Remove(c *gin.Context) {
	uidStr := c.Param("id")
	uid, err := strconv.ParseInt(uidStr, 10, 64)
	if err == nil {
		ui, _ := c.Get("user")
		user := ui.(*models.ServerUser)
		user.RemoveChatUser(uid)
		if uClient, exist := hub.Hub.User.AcceptedClient.GetClient(uid); exist {
			if uClient.ServerId == user.ID {
				hub.Hub.User.Change2waiting(uClient)
			}
		}
	}
	c.JSON(200, util.RespSuccess(nil))
}
