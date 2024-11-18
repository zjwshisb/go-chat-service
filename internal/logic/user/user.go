package user

import (
	"gf-chat/internal/dao"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"
	"gf-chat/internal/trait"
)

func init() {
	service.RegisterUser(&sUser{
		trait.Curd[entity.Users]{
			Dao: &dao.Users,
		},
	})
}

type sUser struct {
	trait.Curd[entity.Users]
}

func (s sUser) FindByToken(token string) *entity.Users {
	return nil
}
