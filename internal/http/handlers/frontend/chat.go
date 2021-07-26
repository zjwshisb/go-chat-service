package frontend

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"ws/configs"
	"ws/internal/auth"
	"ws/internal/chat"
	"ws/internal/file"
	"ws/internal/json"
	"ws/internal/repositories"
	"ws/internal/util"
)

// 消息记录
func GetHistoryMessage(c *gin.Context) {
	user := auth.GetUser(c)
	wheres := []*repositories.Where{
		{
			Filed: "user_id = ?",
			Value: user.GetPrimaryKey(),
		},
	}
	id, exist := c.GetQuery("id")
	if exist {
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err == nil {
			wheres = append(wheres, &repositories.Where{
				Filed: "id < ?",
				Value: idInt,
			})
		}
	}
	messages := repositories.GetMessages(wheres, 100, []string{"BackendUser"})
	messagesResources := make([]*json.Message, 0, len(messages))
	for _, m := range messages {
		var avatar string
		switch m.Source {
		case 0:
			avatar = user.GetAvatarUrl()
		case 1:
			avatar = m.BackendUser.GetAvatarUrl()
		case 2:
			avatar = util.SystemAvatar()
		}
		messagesResources = append(messagesResources, &json.Message{
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
		})
	}
	util.RespSuccess(c, messagesResources)
}
// 获取微信订阅消息ID，只有当前没有订阅的时候才会返回
func GetTemplateId(c *gin.Context) {
	user := auth.GetUser(c)
	id := ""
	if !chat.IsSubScribe(user.GetPrimaryKey()) {
		id = configs.Wechat.SubscribeTemplateIdOne
	}
	util.RespSuccess(c , gin.H{
		"id": id,
	})
}
// 标记已订阅微信订阅消息
func Subscribe(c *gin.Context) {
	user := auth.GetUser(c)
	_ = chat.SetSubscribe(user.GetPrimaryKey())
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
