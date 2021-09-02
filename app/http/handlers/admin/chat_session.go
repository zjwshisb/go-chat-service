package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon"
	"ws/app/databases"
	"ws/app/models"
	"ws/app/repositories"
	"ws/app/util"
)

type ChatSessionHandler struct {

}

// 获取会话列表
func (handler *ChatSessionHandler) Index(c *gin.Context)  {
	sessions := make([]*models.ChatSession, 0 ,0)
	var total int64
	wheres := make([]*repositories.Where, 0)
	if c.Query("admin_name") != "" {
		admins := make([]*models.Admin, 0)
		databases.Db.Where("username like ?" ,
			"%" + c.Query("admin_name")  + "%").
			Find(&admins)
		ids := make([]int64, 0, 0)
		for _, admin := range admins {
			ids = append(ids, admin.ID)
		}
		wheres = append(wheres, &repositories.Where{
			Filed: "admin_id in ?",
			Value: ids,
		})
	}
	queriedAtArr := c.QueryArray("queried_at")
	if len(queriedAtArr) > 0 {
		start := carbon.Parse(queriedAtArr[0]).ToTimestamp()
		wheres = append(wheres, &repositories.Where{
			Filed: "queried_at >= ?",
			Value: start,
		})
		if len(queriedAtArr) > 1 {
			end := carbon.Parse(queriedAtArr[1]).ToTimestamp()
			wheres = append(wheres, &repositories.Where{
				Filed: "queried_at <= ?",
				Value: end,
			})
		}
	}
	databases.Db.Scopes(repositories.Paginate(c)).
		Scopes(repositories.AddWhere(wheres)).
		Preload("Admin").
		Preload("User").
		Order("id desc").
		Find(&sessions)
	databases.Db.Model(&models.ChatSession{}).
		Scopes(repositories.AddWhere(wheres)).
		Count(&total)
	data := make([]*models.ChatSessionJson, 0, len(sessions))
	for _, session := range sessions {
		data = append(data, session.ToJson())
	}
	util.RespPagination(c, repositories.NewPagination(data, total))
}
// 会话详情
func (handler *ChatSessionHandler) Show(c *gin.Context) {
	sessionId := c.Param("id")
	session := &models.ChatSession{}
	databases.Db.Preload("Admin").Preload("User").
		Find(session ,sessionId)
	messages := make([]*models.Message, 0, 0)
	databases.Db.Preload("User").Preload("Admin").
		Where("session_id = ?" , sessionId).
		Where("source in ?", []int{models.SourceAdmin, models.SourceUser}).
		Find(&messages)
	data := make([]*models.MessageJson, 0, 0)
	for _, msg:= range messages {
		data = append(data, msg.ToJson())
	}
	util.RespSuccess(c, gin.H{
		"messages": data,
		"total": len(data),
		"session": session.ToJson(),
	})
}