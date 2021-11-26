package admin

import (
	"github.com/gin-gonic/gin"
	"ws/app/chat"
	"ws/app/http/requests"
	"ws/app/models"
	"ws/app/repositories"
	"ws/app/util"
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
	})
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
	// todo
	//_ , exist := websocket.AdminManager.GetConn(transfer.)
	//if exist {
	//	websocket.AdminManager.BroadcastUserTransfer(transfer.ToAdminId)
	//}
	util.RespSuccess(c , gin.H{})
}

func (handler *TransferHandler) Index(c *gin.Context)  {
	wheres := requests.GetFilterWhere(c, map[string]interface{}{})
	p := transferRepo.Paginate(c, wheres, []string{"User","ToAdmin","FromAdmin"}, "id desc")
	_ = p.DataFormat(func(i interface{}) interface{} {
		item := i.(*models.ChatTransfer)
		return item.ToJson()
	})
	util.RespPagination(c , p)
}
