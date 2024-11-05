package mealorder

import (
	"context"
	"gf-chat/internal/consts"
	"gf-chat/internal/dao"
	"gf-chat/internal/service"
)

func init() {
	service.RegisterMealOrder(&sMealOrder{})
}

type sMealOrder struct {
}

func (s sMealOrder) GetActiveCount(ctx context.Context, uid int, w any) int {
	query := dao.MealOrders.Ctx(ctx).Where("user_id", uid).
		WhereIn("status", []int{consts.MealOrderStatusCommitted, consts.MealOrderStatusMaking})
	if w != nil {
		query = query.Where(w)
	}
	count, _ := query.Count()
	return count
}
