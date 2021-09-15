package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon"
	"time"
	"ws/app/chat"
	"ws/app/http/requests"
	"ws/app/models"
	"ws/app/repositories"
	"ws/app/util"
	"ws/app/websocket"
)

type ChatSessionHandler struct {

}
var sessionFilter = map[string]interface{}{
	"admin_name": func(val string) Where {
		admins := adminRepo.Get([]Where{
			{
				Filed: "username like ?",
				Value: "%" + val + "%",
			},
		}, -1, []string{})
		ids := make([]int64, 0, 0)
		for _, admin := range admins {
			ids = append(ids, admin.ID)
		}
		return &repositories.Where{
			Filed: "admin_id in ?",
			Value: ids,
		}
	},
	"status": func(val string) interface{} {
		switch val {
		case "cancel":
			return &repositories.Where{
				Filed: "canceled_at > ?",
				Value: 0,
			}
		case "accept":
			return &repositories.Where{
				Filed: "accepted_at > ?",
				Value: 0,
			}
		case "wait":
			wheres := []*repositories.Where{
				{
					Filed: "accepted_at = ?",
					Value: 0,
				},
				{
					Filed: "canceled_at = ?",
					Value: 0,
				},
			}
			return wheres
		}
		return nil
	},
}

// 获取会话列表
func (handler *ChatSessionHandler) Index(c *gin.Context)  {
	wheres := requests.GetFilterWhere(c, sessionFilter)
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
		item := i.(*models.ChatSession)
		return item.ToJson()
	})
	util.RespPagination(c, p)
}
func (handler *ChatSessionHandler) Cancel(c *gin.Context)  {
	sessionId := c.Param("id")
	session := chatSessionRepo.First([]Where{
		{
			Filed: "id = ?",
			Value: sessionId,
		},
	})
	if session == nil {
		util.RespNotFound(c)
		return
	}
	if session.AcceptedAt > 0 {
		util.RespFail(c, "会话已接入，无法取消", 500)
		return
	}
	if session.CanceledAt > 0 {
		util.RespFail(c, "会话已取消，请勿重复取消", 500)
		return
	}
	session.CanceledAt = time.Now().Unix()
	chatSessionRepo.Save(session)
	_ = chat.RemoveManual(session.UserId)
	websocket.AdminHub.BroadcastWaitingUser()
	util.RespSuccess(c, gin.H{})
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