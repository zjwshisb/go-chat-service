package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon"
	"strconv"
	"ws/app/chat"
	"ws/app/http/requests"
	"ws/app/models"
	"ws/app/repositories"
	"ws/app/resource"
	"ws/app/util"
	"ws/app/websocket"
)

type AdminsHandler struct {
}

func (handle *AdminsHandler) Index(c *gin.Context){
	where := requests.GetFilterWhere(c, map[string]interface{}{
		"username": "=",
	})
	admin := requests.GetAdmin(c)
	where = append(where, &repositories.Where{
		Filed: "group_id = ?",
		Value: admin.GetGroupId(),
	})
	p := repositories.AdminRepo.Paginate(c, where, []string{}, []string{"id desc"})
	_ = p.DataFormat(func(i interface{}) interface{} {
		admin := i.(*models.Admin)
		return &resource.Admin{
			Avatar:        admin.GetAvatarUrl(),
			Username:      admin.Username,
			Online:        websocket.AdminManager.ConnExist(admin),
			Id:            admin.ID,
			AcceptedCount: chat.AdminService.GetActiveCount(admin.GetPrimaryKey()),
		}
	})
	util.RespPagination(c, p)
}

func (handle *AdminsHandler) Show(c *gin.Context){
	id := c.Param("id")
	u := requests.GetAdmin(c)
	admin := repositories.AdminRepo.First([]*repositories.Where{
		{
			Filed: "id = ?",
			Value: id,
		},
		{
			Filed: "group_id = ?",
			Value: u.GetGroupId(),
		},
	}, []string{})
	if admin == nil {
		util.RespNotFound(c)
		return
	}
	month := c.Query("month")
	date := carbon.Parse(month)
	if date.Error != nil {
		date = carbon.Now()
	}
	firstDate := date.StartOfMonth()
	lastDate := date.EndOfMonth()
	firstDateUnix := firstDate.ToTimestamp()
	sessions := repositories.ChatSessionRepo.Get([]*repositories.Where{
		{
			Filed: "accepted_at >= ?",
			Value: firstDateUnix,
		},
		{
			Filed: "accepted_at <= ?",
			Value: lastDate.ToTimestamp(),
		},
		{
			Filed: "admin_id = ?",
			Value: admin.GetPrimaryKey(),
		},
	}, -1,  []string{}, []string{"id"})
	
	value := make([]*resource.Line, lastDate.DayOfMonth())

	for day, _ := range value {
		value[day] = &resource.Line{
			Category: "每日接待数",
			Value:    0,
			Label:    strconv.Itoa(day + 1) + "号",
		}
	}
	for _, session := range sessions {
		d := (session.AcceptedAt - firstDateUnix) / (24 * 3600)
		value[d].Value += 1
	}
	util.RespSuccess(c, gin.H{
		"chart": value,
		"admin": resource.Admin{
			Avatar:        admin.GetAvatarUrl(),
			Username:      admin.GetUsername(),
			Online:        websocket.AdminManager.IsOnline(admin),
			Id:            admin.GetPrimaryKey(),
			AcceptedCount: chat.AdminService.GetActiveCount(admin.GetPrimaryKey()),
		},
	})
}


