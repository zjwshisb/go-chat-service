package admin

import (
	"ws/app/chat"
	"ws/app/http/requests"
	"ws/app/http/responses"
	"ws/app/http/websocket"
	"ws/app/models"
	"ws/app/repositories"

	"github.com/gin-gonic/gin"
)

type TransferHandler struct {
}

func (handler *TransferHandler) Cancel(c *gin.Context) {
	id := c.Param("id")
	transfer := repositories.TransferRepo.First([]*repositories.Where{
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
		responses.RespNotFound(c)
		return
	}
	if transfer.IsCanceled {
		responses.RespValidateFail(c, "transfer is canceled")
		return
	}
	if transfer.IsAccepted {
		responses.RespValidateFail(c, "transfer is accepted")
		return
	}
	_ = chat.TransferService.Cancel(transfer)
	websocket.AdminManager.NoticeUserTransfer(requests.GetAdmin(c))
	responses.RespSuccess(c, gin.H{})
}

func (handler *TransferHandler) Index(c *gin.Context) {
	wheres := requests.GetFilterWhere(c, map[string]interface{}{})
	wheres = append(wheres, &repositories.Where{
		Filed: "group_id = ?",
		Value: requests.GetAdmin(c).GetGroupId(),
	})
	p := repositories.TransferRepo.Paginate(c, wheres, []string{"User", "ToAdmin", "FromAdmin"}, []string{"id desc"})
	_ = p.DataFormat(func(item *models.ChatTransfer) interface{} {
		return item.ToJson()
	})
	responses.RespPagination(c, p)
}
