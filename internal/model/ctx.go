package model

import (
	"github.com/gogf/gf/v2/frame/g"
)

type AdminCtx struct {
	Entity *CustomerAdmin
	Data   g.Map
}

type UserCtx struct {
	Entity *User
	Data   g.Map
}
