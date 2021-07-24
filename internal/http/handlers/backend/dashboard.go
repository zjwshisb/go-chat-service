package backend

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon"
	"gorm.io/gorm"
	"sort"
	"ws/internal/chat"
	"ws/internal/databases"
	"ws/internal/models"
	"ws/internal/util"
	"ws/internal/websocket"
)

func GetUserQueryInfo(c *gin.Context) {
	startTime := carbon.Now().StartOfDay().ToTimestamp()
	endTime := carbon.Now().EndOfDay().ToTimestamp()
	records := make([]models.ChatSession, 0)
	static := make(map[int64]map[string]int64)
	var i int64
	for i = 0; i<=23; i++ {
		item := make(map[string]int64)
		item["count"] = 0
		static[i] = item
	}
	var total int64
	databases.Db.Table("query_records").
		Where("queried_at >= ?", startTime).
		Where("queried_at <= ?", endTime).
		Count(&total)


	var messageCount int64
	databases.Db.Table("messages").
		Where("received_at >= ?", startTime).
		Where("received_at <= ?" , endTime).
		Where("source = ?" , models.SourceUser).
		Count(&messageCount)

	var totalTime int64
	var maxTime int64
	var acceptCount int64

	databases.Db.Table("query_records").
		Order("queried_at desc").
		Where("queried_at >= ?", startTime).
		Where("queried_at <= ?", endTime).
		FindInBatches(&records,
			100,
			func(tx *gorm.DB, batch int) error {
				for _, model := range records {
					hour := (model.QueriedAt - startTime) / 3600
					item ,exist := static[hour]
					if exist {
						item["count"] = item["count"] + 1
					}
					if model.ServiceId > 0 {
						dura :=  model.AcceptedAt - model.QueriedAt
						totalTime += dura
						if dura > maxTime {
							maxTime = dura
						}
						acceptCount += 1
					}
				}
				return nil
			})
	resp := make([]map[string]interface{}, 0)
	for hour, i := range static{
		userItem := make(map[string]interface{})
		userItem["category"] = "用户数"
		userItem["label"] = hour
		userItem["count"] = i["count"]
		resp = append(resp, userItem)
	}
	sort.Slice(resp, func(i, j int) bool {
		return resp[i]["label"].(int64) < resp[j]["label"].(int64)
	})
	util.RespSuccess(c , gin.H{
		"user_count" : total,
		"message_count": messageCount,
		"avg_time": totalTime / acceptCount,
		"max_time" : maxTime,
		"chart": resp,
	})
}

func GetOnlineInfo(c *gin.Context)  {
	util.RespSuccess(c, gin.H{
		"user_count": len(websocket.UserHub.GetAllConn()),
		"service_count": len(websocket.ServiceHub.GetAllConn()),
		"waiting_user_count": len(chat.GetManualUserIds()),
	})
}
