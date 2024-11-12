// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"github.com/gogf/gf/v2/os/gctx"
)

type (
	ISetup interface {
		Setup(ctx gctx.Ctx, customerId uint)
	}
)

var (
	localSetup ISetup
)

func Setup() ISetup {
	if localSetup == nil {
		panic("implement not found for interface ISetup, forgot register?")
	}
	return localSetup
}

func RegisterSetup(i ISetup) {
	localSetup = i
}
