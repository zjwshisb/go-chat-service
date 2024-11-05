package model

import (
	"gf-chat/internal/model/entity"
	"github.com/gogf/gf/v2/frame/g"
)

type UserCtx struct {
	Entity  *entity.Users
	UserApp *entity.UserApps
	Data    g.Map
}
