package chat

import (
	"gf-chat/internal/contract"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/model/relation"
	"github.com/gogf/gf/v2/util/gconv"
)

type chatConnMessage struct {
	Conn iWsConn
	Msg  *relation.CustomerChatMessages
}

type user struct {
	Entity *entity.Users
}

func (u user) GetUser() any {
	return u.Entity
}

func (u user) GetPrimaryKey() int {
	return u.Entity.Id
}

func (u user) GetCustomerId() int {
	return u.Entity.CustomerId
}

func (u user) GetAvatarUrl() string {
	return ""
}

func (u user) GetUsername() string {
	return u.Entity.Phone
}

func (u user) AccessTo(user contract.IChatUser) bool {
	return true
}

type admin struct {
	Entity *relation.CustomerAdmins
}

func (u admin) GetUser() any {
	return u.Entity
}

func (u admin) GetPrimaryKey() int {
	return gconv.Int(u.Entity.Id)
}

func (u admin) GetCustomerId() int {
	return u.Entity.CustomerId
}

func (u admin) GetAvatarUrl() string {
	return ""
}

func (u admin) GetUsername() string {
	return u.Entity.Username
}

func (u admin) AccessTo(user contract.IChatUser) bool {
	return true
}
