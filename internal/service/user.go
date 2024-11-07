// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
)

type (
	IUser interface {
		GetUsers(ctx context.Context, w any) []*entity.Users
		First(w do.Users) *entity.Users
		FindByToken(token string) *entity.Users
	}
)

var (
	localUser IUser
)

func User() IUser {
	if localUser == nil {
		panic("implement not found for interface IUser, forgot register?")
	}
	return localUser
}

func RegisterUser(i IUser) {
	localUser = i
}
