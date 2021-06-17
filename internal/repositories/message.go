package repositories

import (
	"ws/internal/databases"
	"ws/internal/models"
)
type Where struct {
	Filed string
	Value interface{}
}
func UpdateMessages(wheres []Where, values map[string]interface{}) int64 {
	query := databases.Db.Table("messages")
	for _, where := range wheres {
		query = query.Where(where.Filed, where.Value)
	}
	query.Updates(values)
	return query.RowsAffected
}

func GetMessages(wheres []Where, limit int, load []string) []*models.Message  {
	messages := make([]*models.Message, 0)
	query := databases.Db.Order("received_at desc").Order("id desc").Limit(limit)
	for _, where := range wheres {
		query = query.Where(where.Filed, where.Value)
	}
	for _, relate := range load {
		query = query.Preload(relate)
	}
	query.Find(&messages)
	return messages
}

func GetUnSendMessage(wheres []Where, load []string) []*models.Message {
	wheres = append(wheres, Where{
		Filed: "service_id = ?",
		Value: 0,
	})
	return GetMessages(wheres, -1, load)
}
