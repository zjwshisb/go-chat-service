// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	ISso interface {
		Check(ctx context.Context, sessionId string, uid int) bool
		Auth(ctx context.Context, ticket string) (uid int, sessionId string, err error)
	}
)

var (
	localSso ISso
)

func Sso() ISso {
	if localSso == nil {
		panic("implement not found for interface ISso, forgot register?")
	}
	return localSso
}

func RegisterSso(i ISso) {
	localSso = i
}
