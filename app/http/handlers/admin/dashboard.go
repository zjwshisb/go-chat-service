package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon"
	"gorm.io/gorm"
	"sort"
	"ws/app/chat"
	"ws/app/databases"
	"ws/app/http/requests"
	"ws/app/models"
	"ws/app/util"
	"ws/app/websocket"
)

type DashboardHandler struct {

}

func (handler *DashboardHandler) GetUserQueryInfo(c *gin.Context) {
	startTime := carbon.Now().StartOfDay().ToTimestamp()
	endTime := carbon.Now().EndOfDay().ToTimestamp()
	sessions := make([]models.ChatSession, 0)
	static := make(map[int64]map[string]int64)
	var i int64
	for i = 0; i<=23; i++ {
		item := make(map[string]int64)
		item["count"] = 0
		static[i] = item
	}
	var total int64
	databases.Db.Model(&models.ChatSession{}).
		Where("queried_at >= ?", startTime).
		Where("queried_at <= ?", endTime).
		Count(&total)
	var messageCount int64
	databases.Db.Model(&models.Message{}).
		Where("received_at >= ?", startTime).
		Where("received_at <= ?" , endTime).
		Where("source = ?" , models.SourceUser).
		Count(&messageCount)

	var totalTime int64
	var maxTime int64
	var acceptCount int64

	databases.Db.
		Order("queried_at desc").
		Where("queried_at >= ?", startTime).
		Where("queried_at <= ?", endTime).
		Where("accepted_at > ?", 0).
		FindInBatches(&sessions,
			100,
			func(tx *gorm.DB, batch int) error {
				for _, model := range sessions {
					hour := (model.QueriedAt - startTime) / 3600
					item ,exist := static[hour]
					if exist {
						item["count"] = item["count"] + 1
					}
					if model.AdminId > 0 {
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
	var avgTime int64
	if acceptCount > 0 {
		avgTime = totalTime / acceptCount
	}
	util.RespSuccess(c , gin.H{
		"user_count" : total,
		"message_count": messageCount,
		"avg_time": avgTime,
		"max_time" : maxTime,
		"chart": resp,
	})
}

func (handler *DashboardHandler) GetOnlineInfo(c *gin.Context)  {
	admin := requests.GetAdmin(c)
	util.RespSuccess(c, gin.H{
		"user_count": websocket.UserManager.GetOnlineTotal(admin.GetGroupId()),
		"admin_count": websocket.AdminManager.GetOnlineTotal(admin.GetGroupId()),
		"waiting_user_count": len(chat.ManualService.GetAll(admin.GetGroupId())),
	})
}
