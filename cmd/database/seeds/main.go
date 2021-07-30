package main

import (
	"gorm.io/gorm"
	"ws/app/databases"
	"ws/app/models"
)

func init()  {
	databases.Setup()
}
func rules()  {
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
			Match: models.MatchServiceAllOffLine,
			ReplyType: models.ReplyTypeMessage,
			IsSystem: 1,
		},
	}
	for _, rule := range rules {
		query := databases.Db.
			Where("`match`=?" , rule.Match).
			First(&models.AutoRule{})
		if query.Error == gorm.ErrRecordNotFound {
			databases.Db.Save(&rule)
		}
	}
}

func main()  {
	rules()
}

