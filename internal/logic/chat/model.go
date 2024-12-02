package chat

import (
	"gf-chat/internal/model"
)

type IChatUser interface {
	GetUser() any
	GetPrimaryKey() uint
	GetUsername() string
	GetAvatarUrl() string
	GetCustomerId() uint
	AccessTo(user IChatUser) bool
}

type chatConnMessage struct {
	Conn iWsConn
	Msg  *model.CustomerChatMessage
}

type user struct {
	Entity *model.User
}

func (u user) GetUser() any {
	return u.Entity
}

func (u user) GetPrimaryKey() uint {
	return u.Entity.Id
}

func (u user) GetCustomerId() uint {
	return u.Entity.CustomerId
}

func (u user) GetAvatarUrl() string {
	return ""
}

func (u user) GetUsername() string {
	return u.Entity.Username
}

func (u user) AccessTo(user IChatUser) bool {
	return true
}

type admin struct {
	Entity *model.CustomerAdmin
}

func (u admin) GetUser() any {
	return u.Entity
}

func (u admin) GetPrimaryKey() uint {
	return u.Entity.Id
}

func (u admin) GetCustomerId() uint {
	return u.Entity.CustomerId
}

func (u admin) GetAvatarUrl() string {
	return ""
}

func (u admin) GetUsername() string {
	return u.Entity.Username
}

func (u admin) AccessTo(user IChatUser) bool {
	return true
}
