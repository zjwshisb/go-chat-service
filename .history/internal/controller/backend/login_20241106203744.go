package backend

import (
	"context"
	"database/sql"
	api "gf-chat/api/v1/backend/user"
	"gf-chat/internal/dao"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
	"golang.org/x/crypto/bcrypt"
)

var CLogin = &cLogin{}

type cLogin struct {
}

func (login *cLogin) Login(ctx context.Context, r *api.LoginReq) (res *api.LoginRes, err error) {
	admin := &entity.CustomerAdmins{}
	err = dao.CustomerAdmins.Ctx(ctx).Where(do.CustomerAdmins{Username: r.Username}).Scan(admin)
	if err == sql.ErrNoRows {
		return nil, gerror.NewCode(gcode.CodeValidationFailed, "账号或密码错误")
	}
	err = service.Admin().IsValid(admin)
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, err.Error())
	}
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(r.Password))
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeValidationFailed, "账号或密码错误")
	}
	token, _ := service.Jwt().CreateToken(gconv.String(admin.Id), "")
	return &api.LoginRes{Token: token}, nil
}
