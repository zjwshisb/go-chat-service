package refund

import (
	"context"
	"gf-chat/internal/consts"
	"gf-chat/internal/dao"
	"gf-chat/internal/service"
)

func init() {
	service.RegisterRefund(&sRefund{})
}

type sRefund struct {
}

func (s sRefund) GetCommittedCount(ctx context.Context, uid int) int {
	count, _ := dao.RefundOrders.Ctx(ctx).
		Where("user_id", uid).
		Where("status", consts.RefundOrderStatusCommitted).Count()
	return count
}
