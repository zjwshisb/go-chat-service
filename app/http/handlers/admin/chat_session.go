package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon"
	"ws/app/models"
	"ws/app/repositories"
	"ws/app/util"
)

type ChatSessionHandler struct {

}

// 获取会话列表
func (handler *ChatSessionHandler) Index(c *gin.Context)  {
	wheres := make([]*repositories.Where, 0)
	if c.Query("admin_name") != "" {
		admins := adminRepo.Get([]Where{
			{
				Filed: "username like ?",
				Value: "%" + c.Query("admin_name")  + "%",
			},
		}, -1 ,[]string{})
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
	p := chatSessionRepo.Paginate(c , wheres, []string{"Admin","User"},"id desc")
	_ = p.DataFormat(func(i interface{}) interface{} {
		item := i.(models.ChatSession)
		return item.ToJson()
	})
	util.RespPagination(c, p)
}
// 会话详情
func (handler *ChatSessionHandler) Show(c *gin.Context) {
	sessionId := c.Param("id")
	session := chatSessionRepo.First([]Where{
		{
			Filed: "id = ?",
			Value: sessionId,
		},
	})
	messages := messageRepo.Get([]Where{
		{
			Filed: "session_id = ?",
			Value: sessionId,
		},
		{
			Filed: "source in ?",
			Value: []int{models.SourceAdmin, models.SourceUser},
		},
	}, -1, []string{"User", "Admin"})
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