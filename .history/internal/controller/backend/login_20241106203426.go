package backend

import (
	"context"
	"database/sql"
	api "gf-chat/api/v1/backend/user"
	"gf-chat/internal/dao"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"
	"strconv"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
	"golang.org/x/crypto/bcrypt"
)

var CLogin = &cLogin{}

type cLogin struct {
}

func (login *cLogin) Password(ctx context.Context, r *api.LoginReq) (res *api.LoginReq, err error) {
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
	token, err := service.Jwt().CreateToken(gconv.String(admin.Id), "")
	return &api.LoginRes{Token: token}, nil
}

func (login *cLogin) Login(ctx context.Context, r *api.LoginSSoReq) (res *api.LoginRes, err error) {
	ticket := r.Ticket
	uid, sessionId, e := service.Sso().Auth(ctx, ticket)
	if e != nil {
		err = gerror.NewCode(gcode.CodeBusinessValidationFailed, e.Error())
		return
	}
	token, _ := service.Jwt().CreateToken(strconv.Itoa(uid), sessionId)
	res = &api.LoginRes{Token: token}
	return
}
