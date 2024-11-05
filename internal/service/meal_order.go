// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package service

import (
	"context"
)

type IMealOrder interface {
	GetActiveCount(ctx context.Context, uid int, w any) int
}

var localMealOrder IMealOrder

func MealOrder() IMealOrder {
	if localMealOrder == nil {
		panic("implement not found for interface IMealOrder, forgot register?")
	}
	return localMealOrder
}

func RegisterMealOrder(i IMealOrder) {
	localMealOrder = i
}
