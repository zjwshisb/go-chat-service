package service

import (
	"gf-chat/internal/model/entity"
	"gf-chat/internal/trait"
)

type (
	IUser interface {
		trait.ICurd[entity.Users]
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
