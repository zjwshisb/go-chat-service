package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon"
	"time"
	"ws/app/chat"
	"ws/app/http/requests"
	"ws/app/models"
	"ws/app/repositories"
	"ws/app/resource"
	"ws/app/util"
	"ws/app/websocket"
)

type ChatSessionHandler struct {

}
var sessionFilter = map[string]interface{}{
	"admin_name": func(val string) *repositories.Where {
		admins := repositories.AdminRepo.Get([]*repositories.Where{
			{
				Filed: "username like ?",
				Value: "%" + val + "%",
			},
		}, -1, []string{}, []string{})
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

// Index 获取会话列表
func (handler *ChatSessionHandler) Index(c *gin.Context)  {
	wheres := requests.GetFilterWhere(c, sessionFilter)
	queriedAtArr := c.QueryArray("queried_at")
	wheres = append(wheres,&repositories.Where{
		Filed: "group_id = ?",
		Value: requests.GetAdmin(c).GetGroupId(),
	})
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
	p := repositories.ChatSessionRepo.Paginate(c , wheres, []string{"Admin","User"},[]string{"id desc"})
	_ = p.DataFormat(func(i interface{}) interface{} {
		item := i.(*models.ChatSession)
		return item.ToJson()
	})
	util.RespPagination(c, p)
}
func (handler *ChatSessionHandler) Cancel(c *gin.Context)  {
	sessionId := c.Param("id")
	session := repositories.ChatSessionRepo.First([]*repositories.Where{
		{
			Filed: "id = ?",
			Value: sessionId,
		},
		{
			Filed: "group_id = ?",
			Value: requests.GetAdmin(c).GetGroupId(),
		},
	}, []string{})
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
	repositories.ChatSessionRepo.Save(session)
	_ = chat.ManualService.Remove(session.UserId, session.GetUser().GetGroupId())
	websocket.AdminManager.PublishWaitingUser(session.GetUser().GetGroupId())
	util.RespSuccess(c, gin.H{})
}

// Show 会话详情
func (handler *ChatSessionHandler) Show(c *gin.Context) {
	sessionId := c.Param("id")
	session := repositories.ChatSessionRepo.First([]*repositories.Where{
		{
			Filed: "id = ?",
			Value: sessionId,
		},
		{
			Filed: "group_id = ?",
			Value: requests.GetAdmin(c).GetGroupId(),
		},
	}, []string{})
	messages := repositories.MessageRepo.Get([]*repositories.Where{
		{
			Filed: "session_id = ?",
			Value: sessionId,
		},
		{
			Filed: "group_id = ?",
			Value: requests.GetAdmin(c).GetGroupId(),
		},
		{
			Filed: "source in ?",
			Value: []int{models.SourceAdmin, models.SourceUser},
		},
	}, -1, []string{"User", "Admin"}, []string{"id desc"})
	data := make([]*resource.Message, 0, 0)
	for _, msg:= range messages {
		data = append(data, msg.ToJson())
	}
	util.RespSuccess(c, gin.H{
		"messages": data,
		"total": len(data),
		"session": session.ToJson(),
	})
}