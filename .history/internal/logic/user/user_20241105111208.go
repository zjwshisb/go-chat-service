package user

import (
	"context"
	"database/sql"
	"gf-chat/internal/dao"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"

	"github.com/gogf/gf/v2/os/gctx"
)

func init() {
	service.RegisterUser(&sUser{})
}

type sUser struct {
}

func (s sUser) GetUsers(ctx context.Context, w any) []*entity.Users {
	user := make([]*entity.Users, 0)
	dao.Users.Ctx(ctx).Where(w).
		Scan(&user)
	return user
}

func (s sUser) First(w do.Users) *entity.Users {
	user := &entity.Users{}
	err := dao.Users.Ctx(gctx.New()).Where(w).Scan(user)
	if err == sql.ErrNoRows {
		return nil
	}
	return user
}

func (s sUser) FindByToken(token string) *entity.Users {
	return nil
}
