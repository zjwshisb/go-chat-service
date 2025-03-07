package chat

import (
	"gf-chat/internal/model"
)

type iChatUser interface {
	getUser() any
	getPrimaryKey() uint
	getUsername() string
	getAvatarUrl() string
	getCustomerId() uint
	accessTo(user iChatUser) bool
}

type user struct {
	Entity *model.User
}

func (u user) getUser() any {
	return u.Entity
}

func (u user) getPrimaryKey() uint {
	return u.Entity.Id
}

func (u user) getCustomerId() uint {
	return u.Entity.CustomerId
}

func (u user) getAvatarUrl() string {
	return ""
}

func (u user) getUsername() string {
	return u.Entity.Username
}

func (u user) accessTo(user iChatUser) bool {
	return true
}

type admin struct {
	Entity *model.CustomerAdmin
}

func (u admin) getUser() any {
	return u.Entity
}

func (u admin) getPrimaryKey() uint {
	return u.Entity.Id
}

func (u admin) getCustomerId() uint {
	return u.Entity.CustomerId
}

func (u admin) getAvatarUrl() string {
	return ""
}

func (u admin) getUsername() string {
	return u.Entity.Username
}

func (u admin) accessTo(user iChatUser) bool {
	return true
}
