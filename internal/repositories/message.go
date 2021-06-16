package repositories

import (
	"ws/internal/databases"
	"ws/internal/models"
)

func GetMessages(where map[string]interface{}) []*models.Message  {
	messages := make([]*models.Message, 0)
	query := databases.Db.Preload("User").
		Order("received_at desc").
		Where("service_id = ?", 0).
		Where(where)
	query.Find(&messages)
	return messages
}
func GetUnSendMessage(where map[string]interface{}) []*models.Message {
	messages := make([]*models.Message, 0)
	query := databases.Db.Preload("User").
		Order("received_at desc").
		Where("service_id = ?", 0).
		Where(where)
	query.Find(&messages)
	return messages
}
