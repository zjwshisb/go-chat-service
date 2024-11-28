package model

import (
	"gf-chat/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
)

type AdminCtx struct {
	Entity *CustomerAdmin
	Data   g.Map
}

type UserCtx struct {
	Entity *entity.Users
	Data   g.Map
}