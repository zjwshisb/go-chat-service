package backend

import (
	"context"
	"database/sql"
	baseApi "gf-chat/api"
	"gf-chat/api/v1/backend/user"
	"gf-chat/internal/dao"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
	"golang.org/x/crypto/bcrypt"
)

var CUser = &cUser{}

type cUser struct {
}

func (c *cUser) Index(ctx context.Context, req *user.InfoReq) (res *baseApi.NormalRes[user.InfoRes], err error) {
	admin := service.AdminCtx().GetAdmin(ctx)
	res = baseApi.NewResp(user.InfoRes{
		Id:         admin.Id,
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
		Username:   admin.Username,
	})
	return
}

func (c *cUser) UpdateSetting(ctx context.Context, req *user.UpdateSettingReq) (res *baseApi.NilRes, err error) {
	setting, err := service.Admin().GetSetting(ctx, service.AdminCtx().GetAdmin(ctx))
	if err != nil {
		return nil, err
	}
	setting.Name = req.Name
	setting.Background = req.Background.Path
	setting.Avatar = req.Avatar.Path
	if req.IsAutoAccept {
		setting.IsAutoAccept = 1
	} else {
		setting.IsAutoAccept = 0
	}
	setting.WelcomeContent = req.WelcomeContent
	setting.OfflineContent = req.OfflineContent
	dao.CustomerAdminChatSettings.Ctx(ctx).Save(setting)
	service.Chat().UpdateAdminSetting(service.AdminCtx().GetCustomerId(ctx), setting)
	return &baseApi.NilRes{}, nil
}

func (c *cUser) GetSetting(ctx context.Context, req *user.SettingReq) (res *baseApi.NormalRes[user.SettingRes], err error) {
	admin := service.AdminCtx().GetAdmin(ctx)
	setting, err := service.Admin().GetSetting(ctx, service.AdminCtx().GetAdmin(ctx))
	if err != nil {
		return nil, err
	}
	//avatar := service.Qiniu().Form(setting.Avatar)
	//if avatar == nil {
	avatar := service.ChatSetting().DefaultAvatarForm(admin.CustomerId)
	//}
	return baseApi.NewResp(user.SettingRes{
		//Background:     service.Qiniu().Form(setting.Background),
		IsAutoAccept:   setting.IsAutoAccept == 1,
		WelcomeContent: setting.WelcomeContent,
		OfflineContent: setting.OfflineContent,
		Name:           setting.Name,
		Avatar:         avatar,
	}), nil
}

func (c *cUser) Login(ctx context.Context, r *user.LoginReq) (res *baseApi.NormalRes[user.LoginRes], err error) {
	admin, err := service.Admin().First(ctx, do.CustomerAdmins{Username: r.Username})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gerror.NewCode(gcode.CodeValidationFailed, "账号或密码错误")
		}
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(r.Password))
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeValidationFailed, "账号或密码错误")
	}
	canAccess := service.Admin().CanAccess(admin)
	if !canAccess {
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "账号已禁用")
	}

	token, _ := service.Jwt().CreateToken(gconv.String(admin.Id), "")
	return baseApi.NewResp(user.LoginRes{
		Token: token,
	}), nil
}
