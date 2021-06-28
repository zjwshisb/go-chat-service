package main

import (
	"gorm.io/gorm"
	"ws/internal/databases"
	"ws/internal/models"
)

func setting()  {
	settings := []models.Setting{
		{
			Id: "auto_transfer",
			Name: "自动转接人工客服",
			Description: "当用户咨询时是否自动转接到人工客服",
			Value: "1",
		},
		{
			Id: "session_time",
			Name: "会话有效期(天)",
			Description: "客服与用户的对话有效期，既没有对话后多久失效",
			Value: "1",
		},
	}
	for _, setting := range settings {
		query := databases.Db.
			Where("id = ? " , setting.Id).
			First(&models.Setting{})
		if query.Error == gorm.ErrRecordNotFound {
			databases.Db.Save(&setting)
		}
	}
}