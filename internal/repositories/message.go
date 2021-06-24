package repositories

import (
	"ws/internal/databases"
	"ws/internal/models"
)


func UpdateMessages(wheres []*Where, values map[string]interface{}) int64 {
	query := databases.Db.Table("messages").Scopes(Filter(wheres))
	query.Updates(values)
	return query.RowsAffected
}

func GetMessages(wheres []*Where, limit int, load []string) []*models.Message  {
	messages := make([]*models.Message, 0)
	query := databases.Db.Order("received_at desc").
		Order("id desc").
		Limit(limit).
		Scopes(Filter(wheres))
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

func GetAutoMessage(wheres []*Where, page int, limit int) *Pagination {
	messages := make([]*models.AutoMessage, 0)
	databases.Db.Order("id desc").
		Scopes(Paginate(page, limit)).
		Scopes(Filter(wheres)).
		Limit(limit).
		Find(&messages)
	var total int64
	databases.Db.Model(&models.AutoMessage{}).
		Scopes(Filter(wheres)).
		Count(&total)
	return NewPagination(messages, total)
}