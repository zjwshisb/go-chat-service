package backend

import (
	"github.com/gin-gonic/gin"
	"ws/internal/databases"
	"ws/internal/models"
	"ws/internal/util"
)

func GetChatSession(c *gin.Context)  {
	sessions := make([]*models.ChatSession, 0 ,0)
	var total int64
	databases.Db.Scopes(databases.Paginate(c)).
		Preload("BackendUser").Preload("User").Order("id desc").Find(&sessions)
	databases.Db.Model(&models.ChatSession{}).Order("id desc").Count(&total)
	 data := make([]*models.ChatSessionJson, 0, len(sessions))
	for _, session := range sessions {
		data = append(data, session.ToJson())
	}
	util.RespPagination(c, databases.NewPagination(data, total))
}
func GetChatSessionDetail(c *gin.Context) {
	sessionId := c.Param("id")
	session := &models.ChatSession{}
	databases.Db.Preload("BackendUser").Preload("User").Find(session ,sessionId)
	messages := make([]*models.Message, 0, 0)
	databases.Db.Preload("User").Preload("BackendUser").
		Where("session_id = ?" , sessionId).
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