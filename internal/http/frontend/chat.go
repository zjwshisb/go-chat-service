package frontend

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"ws/configs"
	"ws/internal/auth"
	"ws/internal/file"
	"ws/internal/json"
	"ws/internal/repositories"
	"ws/util"
)

// 消息记录
func GetHistoryMessage(c *gin.Context) {
	user := auth.GetUser(c)
	wheres := []repositories.Where{
		{
			Filed: "user_id = ?",
			Value: user.GetPrimaryKey(),
		},
	}
	id, exist := c.GetQuery("id")
	if exist {
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err == nil {
			wheres = append(wheres, repositories.Where{
				Filed: "id < ?",
				Value: idInt,
			})
		}
	}
	messages := repositories.GetMessages(wheres, 100, []string{"BackendUser"})
	messagesResources := make([]*json.Message, 0, len(messages))
	for _, m := range messages {
		messagesResources = append(messagesResources, &json.Message{
			Id:         m.Id,
			UserId:     m.UserId,
			ServiceId:  m.ServiceId,
			Type:       m.Type,
			Content:    m.Content,
			ReceivedAT: m.ReceivedAT,
			IsServer:   m.IsServer,
			ReqId:      m.ReqId,
			IsRead:     m.IsRead,
			Avatar:     m.BackendUser.GetAvatarUrl(),
		})
	}
	util.RespSuccess(c, messagesResources)
}

func GetTemplateId(c *gin.Context) {
	util.RespSuccess(c , gin.H{
		"id": configs.Wechat.SubscribeTemplateIdOne,
	})
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
