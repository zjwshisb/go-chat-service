package user

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"strconv"
	"ws/app/chat"
	"ws/app/file"
	"ws/app/http/requests"
	"ws/app/repositories"
	"ws/app/resource"
)

// GetHistoryMessage 消息记录
func GetHistoryMessage(c *gin.Context) {
	user := requests.GetUser(c)
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
	var size = 30
	sizeStr, exist := c.GetQuery("size")
	if exist {
		sizeInt, err := strconv.Atoi(sizeStr)
		if err == nil {
			size = sizeInt
		}
	}
	messages := repositories.MessageRepo.Get(wheres, size, []string{"Admin", "User"}, []string{"id desc"})
	messagesResources := make([]*resource.Message, 0, len(messages))
	for _, m := range messages {
		messagesResources = append(messagesResources, m.ToJson())
	}
	requests.RespSuccess(c, messagesResources)
}
func GetReqId(c *gin.Context) {
	u := requests.GetUser(c)
	requests.RespSuccess(c, gin.H{
		"reqId": u.GetReqId(),
	})
}

// GetTemplateId 获取微信订阅消息ID，只有当前没有订阅的时候才会返回
func GetTemplateId(c *gin.Context) {
	user := requests.GetUser(c)
	id := ""
	if !chat.SubScribeService.IsSet(user.GetPrimaryKey()) {
		id = viper.GetString("config.Wechat.SubscribeTemplateIdOne")
	}
	requests.RespSuccess(c, gin.H{
		"id": id,
	})
}

// Subscribe 标记已订阅微信订阅消息
func Subscribe(c *gin.Context) {
	user := requests.GetUser(c)
	_ = chat.SubScribeService.Set(user.GetPrimaryKey())
	requests.RespSuccess(c, gin.H{})
}

// Image 聊天图片
func Image(c *gin.Context) {
	f, _ := c.FormFile("file")
	ff, err := file.Save(f, "chat")
	if err != nil {
		requests.RespFail(c, err.Error(), 500)
	} else {
		requests.RespSuccess(c, gin.H{
			"url": ff.FullUrl,
		})
	}
}
