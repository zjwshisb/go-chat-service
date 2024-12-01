package user

import (
	"context"
	"gf-chat/internal/dao"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"
	"gf-chat/internal/trait"
	"gf-chat/internal/util"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	service.RegisterUser(&sUser{
		trait.Curd[model.User]{
			Dao: &dao.Users,
		},
	})
}

type sUser struct {
	trait.Curd[model.User]
}

func (s *sUser) Login(ctx context.Context, request *ghttp.Request) (user *model.User, token string, err error) {
	username := request.Get("username")
	password := request.Get("password")
	user, err = s.First(ctx, do.Users{Username: username.String()})
	if err != nil {
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), password.Bytes())
	if err != nil {
		err = gerror.NewCode(gcode.CodeValidationFailed, "账号或密码错误")
		return
	}

	token, err = service.Jwt().CreateToken(gconv.String(user.Id))
	if err != nil {
		return
	}
	return
}
func (s *sUser) Auth(ctx g.Ctx, req *ghttp.Request) (user *model.User, err error) {
	token := util.GetRequestToken(req)
	if token == "" {
		err = gerror.NewCode(gcode.CodeNotAuthorized)
		return
	}
	uidStr, err := service.Jwt().ParseToken(token)
	if err != nil {
		err = gerror.NewCode(gcode.CodeNotAuthorized)
		return
	}
	uid := gconv.Int(uidStr)
	user, err = s.Find(ctx, uid)
	if err != nil {
		err = gerror.NewCode(gcode.CodeNotAuthorized)
		return
	}

	return
}
