package main

import (
	"ws/app/databases"
	"ws/app/models"
)

func main()  {
	_ = databases.Db.Migrator().CreateTable(&models.ChatSession{})
	_ = databases.Db.Migrator().CreateTable(&models.Message{})
	_ = databases.Db.Migrator().CreateTable(&models.AutoMessage{})
	_ = databases.Db.Migrator().CreateTable(&models.AdminChatSetting{})
	_ = databases.Db.Migrator().CreateTable(&models.ChatTransfer{})
	_ = databases.Db.Migrator().CreateTable(&models.AutoRule{})
	_ = databases.Db.Migrator().CreateTable(&models.AutoRuleScene{})
	_ = databases.Db.Migrator().CreateTable(&models.Admin{})
	_ = databases.Db.Migrator().CreateTable(&models.User{})

	rules := []models.AutoRule{
		{
			Name: "用户进入客服系统时",
			MatchType: models.MatchTypeAll,
			Match: models.MatchEnter,
			ReplyType: models.ReplyTypeMessage,
			IsSystem: 1,
		},
		{
			Name: "当转接到人工客服而没有客服在线时(如不设置则继续转接到人工客服)",
			MatchType: models.MatchTypeAll,
			Match: models.MatchAdminAllOffLine,
			ReplyType: models.ReplyTypeMessage,
			IsSystem: 1,
		},
	}
	for _, rule := range rules {
		var exist int64
		databases.Db.
			Where("`match`=?" , rule.Match).Count(&exist)
		if exist == 0 {
			databases.Db.Save(&rule)
		}
	}
}

