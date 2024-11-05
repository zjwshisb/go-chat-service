package model

import (
	"gf-chat/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
)

type AdminCtx struct {
	Entity *entity.CustomerAdmins
	Data   g.Map
}
