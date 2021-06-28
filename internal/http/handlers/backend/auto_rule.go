package backend

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"ws/internal/databases"
	"ws/internal/http/requests"
	"ws/internal/json"
	"ws/internal/models"
	"ws/internal/repositories"
	"ws/internal/util"
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
func GetSystemRules(c *gin.Context)  {
	rules := make([]models.AutoRule, 0)
	databases.Db.Where("is_system", 1).Find(&rules)
	util.RespSuccess(c, rules)
}
func UpdateSystemRules(c *gin.Context) {
	 m := make(map[int]int)
	err := c.ShouldBind(&m)
	fmt.Println(m)
	if err != nil {
		util.RespError(c, err.Error())
	}
	for id, v := range m {
		databases.Db.Model(&models.AutoRule{}).
			Where("id = ?", id).
			Update("message_id", v)
	}
	util.RespSuccess(c, m)
}
func GetAutoRules(c *gin.Context)  {
	pagination := repositories.GetAutoRulePagination(c)
	util.RespPagination(c , pagination)
}
func ShowAutoRule(c *gin.Context) {
	id := c.Param("id")
	rule := models.AutoRule{}
	query := databases.Db.Find(&rule, id)
	if query.RowsAffected > 0 {
		util.RespSuccess(c, rule)
	} else {
		util.RespNotFound(c)
	}
}
func StoreAutoRule(c *gin.Context)  {
	form := requests.AutoRuleForm{}
	err := c.BindJSON(&form)
	if err != nil {
		fmt.Println(err)
		util.RespValidateFail(c, err.Error())
		return
	}
	rule := &models.AutoRule{
		Name: form.Name,
		Match: form.Match,
		MatchType: form.MatchType,
		ReplyType: form.ReplyType,
		Sort: form.Sort,
		IsOpen: form.IsOpen,
	}
	if rule.ReplyType == models.ReplyTypeMessage {
		rule.MessageId = form.MessageId
	}
	databases.Db.Save(rule)
	util.RespSuccess(c, rule)
}
func UpdateAutoRule(c *gin.Context) {
	rule := models.AutoRule{}
	result := databases.Db.Find(&rule,  c.Param("id"))
	if result.Error == gorm.ErrRecordNotFound {
		util.RespNotFound(c)
		return
	}
	form := requests.AutoRuleForm{}
	err := c.ShouldBind(&form)
	if err != nil {
		fmt.Println(err)
		util.RespValidateFail(c, err.Error())
		return
	}
	rule.Name = form.Name
	rule.IsOpen = form.IsOpen
	rule.Match = form.Match
	rule.MatchType = form.MatchType
	rule.ReplyType = form.ReplyType
	if rule.ReplyType == models.ReplyTypeTransfer {
		rule.MessageId = 0
	} else {
		rule.MessageId = form.MessageId
	}
	rule.Sort = form.Sort
	rule.MessageId = form.MessageId
	databases.Db.Save(&rule)
	util.RespSuccess(c, rule)
}
func DeleteAutoRule(c *gin.Context)  {

}