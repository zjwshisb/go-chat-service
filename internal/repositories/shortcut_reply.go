package repositories

import (
	"fmt"
	"ws/internal/databases"
	"ws/internal/models"
)
func GetShortcutReply(wheres []Where) []*models.ShortcutReply {
	replies := make([]*models.ShortcutReply, 0)
	query := databases.Db
	for _, where := range wheres {
		query = query.Where(where.Filed, where.Value)
	}
	query.Find(&replies)
	return replies
}

func StoreShortcutReply(data map[string]interface{}) int64 {
	query := databases.Db.Model(&models.ShortcutReply{}).Create(data)
	return query.RowsAffected
}

func DeleteShortcutReply(wheres []Where) int64 {
	query := databases.Db
	for _, where := range wheres {
		query = query.Where(where.Filed, where.Value)
	}
	query.Delete(&models.ShortcutReply{})
	fmt.Println(query.RowsAffected)
	return query.RowsAffected
}