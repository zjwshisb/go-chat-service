package user

import (
	"strconv"
	"time"
	"ws/app/chat"
	"ws/app/file"
	"ws/app/http/requests"
	"ws/app/http/responses"
	"ws/app/models"
	"ws/app/repositories"
	"ws/app/resource"
	"ws/app/util"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func ReadAll(c *gin.Context) {
	form := &struct {
		MsgId int64 `json:"msg_id" binding:"-"`
	}{}
	err := c.Bind(form)
	if err != nil {
		responses.RespValidateFail(c, err)
		return
	}
	wheres := []*repositories.Where{
		{
			Filed: "id <= ?",
			Value: form.MsgId,
		},
		{
			Filed: "user_id = ?",
			Value: requests.GetUser(c).GetPrimaryKey(),
		},
		{
			Filed: "is_read = ?",
			Value: 0,
		},
	}
	repositories.MessageRepo.Update(wheres, map[string]interface{}{
		"is_read": 1,
	})
	responses.RespSuccess(c, gin.H{})
}

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
	messagesResources := make([]*resource.Message, len(messages), len(messages))
	messageIds := make([]int64, len(messages), len(messages))
	for index, m := range messages {
		messagesResources[index] = m.ToJson()
		messageIds[index] = m.Id
	}
	responses.RespSuccess(c, messagesResources)
	// update unread message
	updateWheres := []*repositories.Where{
		{
			Filed: "send_at = ?",
			Value: 0,
		},
		{
			Filed: "source in ?",
			Value: []int{models.SourceSystem, models.SourceAdmin},
		},
		{
			Filed: "id in ?",
			Value: messageIds,
		},
	}
	repositories.MessageRepo.Update(updateWheres, map[string]interface{}{
		"send_at": time.Now().Unix(),
		"is_read": 1,
	})
}

func GetReqId(c *gin.Context) {
	responses.RespSuccess(c, gin.H{
		"reqId": util.RandomStr(20),
	})
}

// GetTemplateId 获取微信订阅消息ID，只有当前没有订阅的时候才会返回
func GetTemplateId(c *gin.Context) {
	user := requests.GetUser(c)
	id := ""
	if !chat.SubScribeService.IsSet(user.GetPrimaryKey()) {
		id = viper.GetString("config.Wechat.SubscribeTemplateIdOne")
	}
	responses.RespSuccess(c, gin.H{
		"id": id,
	})
}

// Subscribe 标记已订阅微信订阅消息
func Subscribe(c *gin.Context) {
	user := requests.GetUser(c)
	_ = chat.SubScribeService.Set(user.GetPrimaryKey())
	responses.RespSuccess(c, gin.H{})
}

// Image 聊天图片
func Image(c *gin.Context) {
	f, _ := c.FormFile("file")
	ff, err := file.Save(f, "chat")
	if err != nil {
		responses.RespFail(c, err.Error(), 500)
	} else {
		responses.RespSuccess(c, gin.H{
			"url": ff.FullUrl,
		})
	}
}
