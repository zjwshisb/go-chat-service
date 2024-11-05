package model

import (
	"gf-chat/internal/model/entity"
	"gf-chat/internal/model/relation"
	"github.com/gogf/gf/v2/frame/g"
)

type AdminCtx struct {
	Entity   *entity.CustomerAdmins
	Data     g.Map
	kitchens []*relation.CustomerKitchen
	schools  []*entity.Schools
}

func (s *AdminCtx) GetKitchens() []*relation.CustomerKitchen {
	return s.kitchens
}

func (s *AdminCtx) SetKitchens(kitchens []*relation.CustomerKitchen) {
	s.kitchens = kitchens
}

func (s *AdminCtx) GetSchools() []*entity.Schools {
	return s.schools
}

func (s *AdminCtx) SetSchools(schools []*entity.Schools) {
	s.schools = schools
}
