package admin

import (
	"github.com/gin-gonic/gin"
	"ws/app/http/requests"
	"ws/app/models"
	"ws/app/repositories"
	"ws/app/resource"
	"ws/app/util"
)

type AutoRuleHandler struct {

}

func (handle *AutoRuleHandler) MessageOptions(c *gin.Context)  {
	messages := autoMessageRepo.Get([]Where{}, -1, []string{}, []string{})
	options := make([]resource.Options, 0, len(messages))
	for _, message := range messages {
		options = append(options, resource.Options{
			Value: message.ID,
			Label: message.Name + "-" + message.TypeLabel(),
		})
	}
	util.RespSuccess(c, options)
}

// SceneOptions 可选择场景
func (handle *AutoRuleHandler) SceneOptions(c *gin.Context) {
	util.RespSuccess(c , models.ScenesOptions)
}

// EventOptions 可选择的事件
func (handle *AutoRuleHandler) EventOptions(c *gin.Context)  {
	util.RespSuccess(c , models.EventOptions)
}

// Index 获取自定义规则列表
func (handle *AutoRuleHandler) Index(c *gin.Context)  {
	filter := map[string]interface{}{
		"reply_type": "=",
		"name" : func(val string) Where {
			return &repositories.Where{
				Filed: "name like ?",
				Value: "%" + val + "%",
			}
		},
		"scenes": func(val string) Where {
			return &repositories.Where{
				Filed: "id in ?",
				Value: autoRuleRepo.GetWithScenesRuleIds(val),
			}
		},
	}
	wheres := requests.GetFilterWhere(c, filter)
	wheres = append(wheres, &repositories.Where{
		Filed: "is_system = ?",
		Value: 0,
	})
	p := autoRuleRepo.Paginate(c, wheres, []string{"Message", "Scenes"}, []string{"id desc"})
	_ = p.DataFormat(func(i interface{}) interface{} {
		item := i.(*models.AutoRule)
		return item.ToJson()
	})
	util.RespPagination(c , p)
}

// Show 获取自定义规则
func (handle *AutoRuleHandler) Show(c *gin.Context) {
	id := c.Param("id")
	rule := autoRuleRepo.First([]Where{
		{
			Filed: "id = ?",
			Value: id,
		},
	}, []string{})
	if rule != nil {
		util.RespSuccess(c, rule.ToJson())
	} else {
		util.RespNotFound(c)
	}
}

// Store 新增自定义规则
func (handle *AutoRuleHandler) Store (c *gin.Context)  {
	form := requests.AutoRuleForm{}
	err := c.ShouldBind(&form)
	if err != nil {
		util.RespValidateFail(c, err.Error())
		return
	}
	if form.ReplyType == models.ReplyTypeTransfer {
		form.Scenes = []string{
			models.SceneNotAccepted,
		}
	}
	rule := &models.AutoRule{
		Name: form.Name,
		Match: form.Match,
		MatchType: form.MatchType,
		ReplyType: form.ReplyType,
		Sort: form.Sort,
		IsOpen: form.IsOpen,
		Key: form.Key,
	}
	var scenes = make([]*models.AutoRuleScene, 0)
	for _, name := range form.Scenes {
		scenes = append(scenes, &models.AutoRuleScene{
			Name:   name,
		})
	}
	rule.Scenes = scenes
	if rule.ReplyType == models.ReplyTypeMessage  || rule.ReplyType == models.ReplyTypeEvent {
		rule.MessageId = form.MessageId
	}
	autoRuleRepo.Save(rule)
	util.RespSuccess(c, rule.ToJson())
}

// Update 更新自定义规则
func (handle *AutoRuleHandler) Update(c *gin.Context) {
	rule := autoRuleRepo.First([]Where{
		{
			Filed: "is_system = ?",
			Value: 0,
		},
		{
			Filed: "id = ?",
			Value: c.Param("id"),
		},
	}, []string{})
	if rule == nil {
		util.RespNotFound(c)
		return
	}
	form := requests.AutoRuleForm{}
	err := c.ShouldBind(&form)
	if err != nil {
		util.RespValidateFail(c, err.Error())
		return
	}
	autoRuleRepo.DeleteScene(rule)
	if form.ReplyType == models.ReplyTypeTransfer {
		form.Scenes = []string{
			models.SceneNotAccepted,
		}
	}
	rule.Name = form.Name
	rule.IsOpen = form.IsOpen
	rule.Match = form.Match
	rule.MatchType = form.MatchType
	rule.ReplyType = form.ReplyType
	rule.Key = form.Key
	if rule.ReplyType == models.ReplyTypeTransfer {
		rule.MessageId = 0
	} else {
		rule.MessageId = form.MessageId
	}
	var scenes = make([]*models.AutoRuleScene, 0)
	for _, name := range form.Scenes {
		scenes = append(scenes, &models.AutoRuleScene{
			Name:   name,
		})
	}
	rule.Scenes = scenes
	rule.Sort = form.Sort
	rule.MessageId = form.MessageId
	autoRuleRepo.Save(rule)
	util.RespSuccess(c, rule.ToJson())
}

// Delete 删除自定义规则
func (handle *AutoRuleHandler) Delete(c *gin.Context)  {
	id := c.Param("id")
	rule := autoRuleRepo.First([]Where{
		{
			Filed: "is_system = ?",
			Value: 0,
		},
		{
			Filed: "id = ?",
			Value: id,
		},
	}, []string{})
	if rule == nil {
		util.RespNotFound(c)
		return
	}
	autoRuleRepo.Delete(rule)
	util.RespSuccess(c, gin.H{})
}

type SystemRuleHandler struct {


}

// Index 获取系统规则
func (handler *SystemRuleHandler) Index(c *gin.Context)  {
	rules := autoRuleRepo.Get([]Where{
		{
			Filed: "is_system = ?",
			Value: 1,
		},
	}, -1, []string{}, []string{})
	result := make([]*resource.AutoRule, len(rules), len(rules))
	for i, rule := range rules {
		result[i] = rule.ToJson()
	}
	util.RespSuccess(c, result)
}

// Update 更新系统规则
func (handler *SystemRuleHandler) Update(c *gin.Context) {
	m := make(map[int]int)
	err := c.ShouldBind(&m)
	if err != nil {
		util.RespError(c, err.Error())
	}
	autoRuleRepo.Update([]Where{
		{
			Filed: "is_system = ?",
			Value: 1,
		},
	}, map[string]interface{}{
		"message_id": 0,
	})
	for id, v := range m {
		autoRuleRepo.Update([]Where{
			{
				Filed: "is_system = ?",
				Value: 1,
			},
			{
				Filed: "id = ?",
				Value: id,
			},
		}, map[string]interface{}{
			"message_id": v,
		})
	}
	util.RespSuccess(c, m)
}