package backend

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon"
	"gorm.io/gorm"
	"ws/internal/databases"
	"ws/internal/models"
	"ws/internal/util"
)

func GetUserQueryInfo(c *gin.Context) {
	startTime := carbon.Now().StartOfDay().ToTimestamp()
	fmt.Println(startTime)
	endTime := carbon.Now().EndOfDay().ToTimestamp()
	records := make([]models.QueryRecord, 0)
	uids := make(map[int64]struct{})
	messageCount := 0
	resp := make(map[int64]map[string]interface{})
	var i int64
	for i = 0; i<=23; i++ {
		item := make(map[string]interface{})
		item["uids"] = make(map[int64]struct{})
		item["message_count"] = 0
		resp[i] = item
	}
	databases.Db.Table("query_records").
		Order("queried_at desc").
		Where("queried_at >= ?", startTime).
		Where("queried_at <= ?", endTime).
		FindInBatches(&records,
			100,
			func(tx *gorm.DB, batch int) error {
				for _, model := range records {
					uids[model.UserId] = struct{}{}
					messageCount += 1
					hour := (model.QueriedAt - startTime) / 3600
					item ,exist := resp[hour]
					if exist {
						item["message_count"] = item["message_count"].(int) + 1
						if u, ok := item["uids"].(map[int64]struct{});ok{
							u[model.UserId] = struct {}{}
						}
					}
				}
				return nil
			})
	for _, item := range resp{
		if m ,ok := item["uids"].(map[int64]struct{}); ok {
			item["user_count"] = len(m)
		}
		delete(item, "uids")
	}
	util.RespSuccess(c , gin.H{
		"user_count" : len(uids),
		"message_count": messageCount,
		"chart": resp,
	})
}
