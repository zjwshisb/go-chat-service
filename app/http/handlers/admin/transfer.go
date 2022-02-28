package admin

import (
	"github.com/gin-gonic/gin"
	"ws/app/chat"
	"ws/app/http/requests"
	"ws/app/models"
	"ws/app/repositories"
	"ws/app/util"
	"ws/app/websocket"
)

type TransferHandler struct {
}


func (handler *TransferHandler) Cancel(c *gin.Context)  {
	id := c.Param("id")
	transfer := transferRepo.First([]*repositories.Where{
		{
			Filed: "id = ?",
			Value: id,
		},
		{
			Filed: "group_id = ?",
			Value: requests.GetAdmin(c).GetGroupId(),
		},
	}, []string{})
	if transfer == nil {
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
	_ = chat.TransferService.Cancel(transfer)
	websocket.AdminManager.PublishTransfer(requests.GetAdmin(c))
	util.RespSuccess(c , gin.H{})
}

func (handler *TransferHandler) Index(c *gin.Context)  {
	wheres := requests.GetFilterWhere(c, map[string]interface{}{})
	wheres = append(wheres, &repositories.Where{
		Filed: "group_id = ?",
		Value: requests.GetAdmin(c).GetGroupId(),
	})
	p := transferRepo.Paginate(c, wheres, []string{"User","ToAdmin","FromAdmin"}, []string{"id desc"})
	_ = p.DataFormat(func(i interface{}) interface{} {
		item := i.(*models.ChatTransfer)
		return item.ToJson()
	})
	util.RespPagination(c , p)
}
