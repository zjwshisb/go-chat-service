package admin

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
	"ws/app/chat"
	"ws/app/databases"
	"ws/app/models"
	"ws/app/repositories"
	"ws/app/util"
	"ws/app/websocket"
)

func CancelTransfer(c *gin.Context)  {
	id := c.Param("id")
	transfer := &models.ChatTransfer{}
	query := databases.Db.Find(transfer, id)
	if query.Error == gorm.ErrRecordNotFound {
		util.RespNotFound(c)
		return
	}
	if transfer.IsCanceled {
		util.RespValidateFail(c, "transfer is canceled")
		return
	}
	if transfer.IsAccepted {
		util.RespValidateFail(c, "transfer is accepted")
		return
	}
	t := time.Now()
	transfer.CanceledAt = &t
	transfer.IsCanceled = true
	_ = chat.CancelTransfer(transfer)
	_ , exist := websocket.AdminHub.GetConn(transfer.ToAdminId)
	if exist {
		websocket.AdminHub.BroadcastUserTransfer(transfer.ToAdminId)
	}
	util.RespSuccess(c , gin.H{})
}

func GetTransfer(c *gin.Context)  {
	transfers := make([]*models.ChatTransfer, 0)
	databases.Db.Order("id desc").
		Scopes(repositories.Paginate(c)).
		Preload("User").Preload("FromAdmin").Preload("ToAdmin").
		Find(&transfers)
	var total int64
	databases.Db.Model(&models.AutoMessage{}).
		Scopes(repositories.Filter(c, []string{"type"})).
		Scopes().
		Count(&total)
	data := make([]*models.ChatTransferJson, 0, len(transfers))
	for _, msg := range transfers {
		data = append(data, msg.ToJson())
	}
	util.RespPagination(c , repositories.NewPagination(data, total))
}
