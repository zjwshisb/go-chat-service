package admin

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"ws/app/databases"
	"ws/app/http/requests"
	"ws/app/json"
	"ws/app/models"
	"ws/app/repositories"
	"ws/app/util"
)

func GetSelectAutoMessage(c *gin.Context)  {
	messages := make([]*models.AutoMessage, 0)
	databases.Db.Find(&messages)
	options := make([]json.Options, 0, len(messages))
	for _, message := range messages {
		options = append(options, json.Options{
			Value: message.ID,
			Label: message.Name + "-" + message.TypeLabel(),
		})
	}
	util.RespSuccess(c, options)
}
func GetSelectScene(c *gin.Context) {
	util.RespSuccess(c , models.ScenesOptions)
}
func GetSelectEvent(c *gin.Context)  {
	util.RespSuccess(c , models.EventOptions)
}

func GetSystemRules(c *gin.Context)  {
	rules := make([]models.AutoRule, 0)
	databases.Db.Where("is_system", 1).Find(&rules)
	result := make([]*models.AutoRuleJson, len(rules), len(rules))
	for i, rule := range rules {
		result[i] = rule.ToJson()
	}
	util.RespSuccess(c, result)
}

func UpdateSystemRules(c *gin.Context) {
	m := make(map[int]int)
	err := c.ShouldBind(&m)
	if err != nil {
		util.RespError(c, err.Error())
	}
	databases.Db.Model(&models.AutoRule{}).
		Where("is_system", 1).
		Update("message_id", 0)
	for id, v := range m {
		databases.Db.Model(&models.AutoRule{}).
			Where("is_system", 1).
			Where("id = ?", id).
			Update("message_id", v)
	}
	util.RespSuccess(c, m)
}
func GetAutoRules(c *gin.Context)  {
	wheres := make([]*repositories.Where, 0)
	name, _ := c.GetQuery("name")
	if name != "" {
		wheres = append(wheres, &repositories.Where{
			Filed: "name like ?",
			Value: "%" + name + "%",
		})
	}
	scene, _ := c.GetQuery("scenes")
	if scene != "" {
		ids := make([]string, 0)
		databases.Db.Model(&models.AutoRuleScene{}).Where("name = ?" , scene).Pluck("rule_id", &ids)
		wheres = append(wheres, &repositories.Where{
			Filed: "id in ?",
			Value: ids,
		})
	}
	pagination := repositories.GetAutoRulePagination(c, wheres)

	util.RespPagination(c , pagination)
}
func ShowAutoRule(c *gin.Context) {
	id := c.Param("id")
	rule := models.AutoRule{}
	query := databases.Db.Preload("Scenes").Where("is_system = ?", 0).Find(&rule, id)
	if query.RowsAffected > 0 {
		util.RespSuccess(c, rule.ToJson())
	} else {
		util.RespNotFound(c)
	}
}
func StoreAutoRule(c *gin.Context)  {
	form := requests.AutoRuleForm{}
	err := c.BindJSON(&form)
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
	databases.Db.Create(rule)
	util.RespSuccess(c, rule.ToJson())
}
func UpdateAutoRule(c *gin.Context) {
	rule := models.AutoRule{}
	result := databases.Db.
		Where("is_system = ?", 0).
		Preload("Scenes").
		Find(&rule,  c.Param("id"))
	for _, s := range rule.Scenes {
		databases.Db.Delete(s)
	}
	if result.Error == gorm.ErrRecordNotFound {
		util.RespNotFound(c)
		return
	}
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
	databases.Db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&rule)
	util.RespSuccess(c, rule.ToJson())
}
func DeleteAutoRule(c *gin.Context)  {
	id := c.Param("id")
	rule := &models.AutoRule{}
	databases.Db.Where("is_system = ?", 0).Find(rule, id)
	if rule.ID == 0 {
		util.RespNotFound(c)
		return
	}
	databases.Db.Delete(rule)
	databases.Db.Table("auto_rule_scenes").Where("rule_id = ?", rule.ID).
		Delete(&models.AutoRuleScene{})
	util.RespSuccess(c, gin.H{})
}