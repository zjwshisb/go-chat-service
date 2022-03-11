package admin

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"ws/app/http/requests"
	"ws/app/models"
	"ws/app/repositories"
	"ws/app/util"
)

type AutoMessageHandler struct {
}

func (handler *AutoMessageHandler) Index(c *gin.Context)  {
	wheres := requests.GetFilterWhere(c, map[string]interface{}{
		"type" : "=",
	})
	admin := requests.GetAdmin(c)
	wheres = append(wheres , &repositories.Where{
		Filed: "group_id = ?",
		Value: admin.GetGroupId(),
	})
	p := repositories.AutoMessageRepo.Paginate(c, wheres, []string{"Rules"}, []string{"id desc"})
	_ = p.DataFormat(func(i interface{}) interface{} {
		item := i.(*models.AutoMessage)
		return item.ToJson()
	})
	util.RespPagination(c , p)
}

func (handler *AutoMessageHandler) Show(c *gin.Context) {
	id := c.Param("id")
	admin := requests.GetAdmin(c)
	message := repositories.AutoMessageRepo.First([]*repositories.Where{
		{
			Filed: "id = ?",
			Value: id,
		},
		{
			Filed: "group_id = ?",
			Value: admin.GetGroupId(),
		},
	}, []string{})
	if message != nil {
		util.RespSuccess(c, message.ToJson())
	} else {
		util.RespNotFound(c)
	}
}

func (handler *AutoMessageHandler) Store(c *gin.Context)  {
	form := requests.AutoMessageForm{}
	err := c.ShouldBind(&form)
	if err != nil {
		util.RespValidateFail(c, err.Error())
		return
	}
	exist := repositories.AutoMessageRepo.First([]*repositories.Where{
		{
			Filed: "name = ?",
			Value: form.Name,
		},
		{
			Filed: "group_id = ?",
			Value: requests.GetAdmin(c).GetGroupId(),
		},
	}, []string{})
	if exist != nil {
		util.RespValidateFail(c, "已存在同名的消息")
		return
	}
	admin := requests.GetAdmin(c)
	message := &models.AutoMessage{
		Name: form.Name,
		Type: form.Type,
		GroupId: admin.GetGroupId(),
	}
	if message.Type == models.TypeText  || message.Type == models.TypeImage {
		message.Content = form.Content
	}
	if message.Type == models.TypeNavigate {
		content := map[string]string{
			"title": form.Title,
			"url": form.Url,
			"content": form.Content,
		}
		jsonBytes, err := json.Marshal(content)
		if err != nil{
			util.RespError(c, err.Error())
			return
		}
		message.Content = string(jsonBytes)
	}
	repositories.AutoMessageRepo.Save(message)
	util.RespSuccess(c, message)
}
func (handler *AutoMessageHandler) Update(c *gin.Context) {
	message := repositories.AutoMessageRepo.First([]*repositories.Where{
		{
			Filed: "id = ?",
			Value: c.Param("id"),
		},
		{
			Filed: "group_id = ?",
			Value: requests.GetAdmin(c).GetGroupId(),
		},
	}, []string{})
	if message == nil {
		util.RespNotFound(c)
		return
	}
	form := requests.AutoMessageForm{}
	err := c.ShouldBind(&form)
	if err != nil {
		util.RespValidateFail(c, err.Error())
		return
	}
	exist := repositories.AutoMessageRepo.First([]*repositories.Where{
		{
			Filed: "name = ?",
			Value: form.Name,
		},
		{
			Filed: "id != ?",
			Value: c.Param("id"),
		},
		{
			Filed: "group_id = ?",
			Value: requests.GetAdmin(c).GetGroupId(),
		},
	}, []string{})
	if exist != nil {
		util.RespValidateFail(c, "已存在同名的其他消息")
		return
	}
	if message.Type == models.TypeText  || message.Type == models.TypeImage {
		message.Content = form.Content
	}
	if message.Type == models.TypeNavigate {
		content := map[string]string{
			"title": form.Title,
			"url": form.Url,
			"content": form.Content,
		}
		jsonBytes, err := json.Marshal(content)
		if err != nil{
			util.RespError(c, err.Error())
			return
		}
		message.Content = string(jsonBytes)
	}
	repositories.AutoMessageRepo.Save(message)
	util.RespSuccess(c, message)
}
func (handler *AutoMessageHandler) Delete(c *gin.Context) {
	message := repositories.AutoMessageRepo.First([]*repositories.Where{
		{
			Filed: "id = ?",
			Value: c.Param("id"),
		},
		{
			Filed: "group_id = ?",
			Value: requests.GetAdmin(c).GetGroupId(),
		},
	}, []string{})
	if message == nil {
		util.RespNotFound(c)
		return
	}
	rules := repositories.AutoRuleRepo.Get([]*repositories.Where{
		{
			Filed: "message_id = ?",
			Value: message.ID,
		},
		{
			Filed: "group_id = ?",
			Value: requests.GetAdmin(c).GetGroupId(),
		},
	}, -1, []string{}, []string{})
	if len(rules) > 0 {
		util.RespValidateFail(c, "该消息在其他地方有使用，无法删除")
		return
	}
	repositories.AutoMessageRepo.Delete(message)
	util.RespSuccess(c, message)
}