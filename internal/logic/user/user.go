package user

import (
	"context"
	"gf-chat/internal/dao"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"
	"gf-chat/internal/trait"
	"github.com/gogf/gf/v2/net/ghttp"
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

func (s sUser) Auth(ctx context.Context, req *ghttp.Request) (*entity.Users, error) {
	return nil, nil
}
