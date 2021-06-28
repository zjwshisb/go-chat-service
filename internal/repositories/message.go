package repositories

import (
	"github.com/gin-gonic/gin"
	"ws/internal/databases"
	"ws/internal/models"
)


func UpdateMessages(wheres []*Where, values map[string]interface{}) int64 {
	query := databases.Db.Table("messages").Scopes(AddWhere(wheres))
	query.Updates(values)
	return query.RowsAffected
}

func GetMessages(wheres []*Where, limit int, load []string) []*models.Message  {
	messages := make([]*models.Message, 0)
	query := databases.Db.Order("received_at desc").
		Order("id desc").
		Limit(limit).
		Scopes(AddWhere(wheres))
	for _, relate := range load {
		query = query.Preload(relate)
	}
	query.Find(&messages)
	return messages
}

func GetUnSendMessage(wheres []*Where, load []string) []*models.Message {
	wheres = append(wheres, &Where{
		Filed: "service_id = ?",
		Value: 0,
	})
	return GetMessages(wheres, -1, load)
}

func GetAutoMessagePagination(c *gin.Context, wheres ...*Where) *databases.Pagination {
	messages := make([]*models.AutoMessage, 0)
	databases.Db.Order("id desc").
		Scopes(databases.Filter(c, []string{"type"})).
		Scopes(databases.Paginate(c)).
		Scopes(AddWhere(wheres)).
		Find(&messages)
	var total int64
	databases.Db.Model(&models.AutoMessage{}).
		Scopes(databases.Filter(c, []string{"type"})).
		Scopes(AddWhere(wheres)).
		Count(&total)
	return databases.NewPagination(messages, total)
}


func GetAutoRulePagination(c *gin.Context, wheres ...*Where) *databases.Pagination {
	rules := make([]*models.AutoRule, 0)
	wheres = append(wheres, &Where{
		Value: 0,
		Filed: "is_system = ?",
	})
	databases.Db.Order("id desc").
		Scopes(databases.Filter(c, []string{"reply_type"})).
		Scopes(databases.Paginate(c)).
		Scopes(AddWhere(wheres)).
		Preload("Message").
		Find(&rules)
	var total int64
	databases.Db.Model(&models.AutoRule{}).
		Scopes(databases.Filter(c, []string{"reply_type"})).
		Scopes(AddWhere(wheres)).
		Count(&total)
	return databases.NewPagination(rules, total)
}