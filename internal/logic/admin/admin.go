package admin

import (
	"context"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/consts"
	"gf-chat/internal/dao"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
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
	service.RegisterAdmin(&sAdmin{
		trait.Curd[model.CustomerAdmin]{
			Dao: &dao.CustomerAdmins,
		},
	})
}

type sAdmin struct {
	trait.Curd[model.CustomerAdmin]
}

func (s *sAdmin) Login(ctx context.Context, request *ghttp.Request) (admin *model.CustomerAdmin, token string, err error) {
	username := request.Get("username")
	password := request.Get("password")
	admin, err = s.First(ctx, do.CustomerAdmins{Username: username.String()})
	if err != nil {
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), password.Bytes())
	if err != nil {
		err = gerror.NewCode(gcode.CodeValidationFailed, "账号或密码错误")
		return
	}
	canAccess := s.CanAccess(admin)
	if !canAccess {
		err = gerror.NewCode(gcode.CodeBusinessValidationFailed, "账号已禁用")
		return
	}
	token, err = service.Jwt().CreateToken(gconv.String(admin.Id))
	if err != nil {
		return
	}
	return
}
func (s *sAdmin) Auth(ctx g.Ctx, req *ghttp.Request) (admin *model.CustomerAdmin, err error) {
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
	admin, err = s.Find(ctx, uid)
	if err != nil {
		err = gerror.NewCode(gcode.CodeNotAuthorized)
		return
	}
	canAccess := s.CanAccess(admin)
	if !canAccess {
		err = gerror.NewCode(gcode.CodeInvalidOperation)
		return
	}
	return
}

func (s *sAdmin) CanAccess(admin *model.CustomerAdmin) bool {
	return true
}

func (s *sAdmin) UpdateSetting(ctx context.Context, admin *model.CustomerAdmin, form api.CurrentAdminSettingForm) (err error) {
	updateData := do.CustomerAdminChatSettings{
		Name:           form.Name,
		IsAutoAccept:   gconv.Int(form.IsAutoAccept),
		WelcomeContent: form.WelcomeContent,
		OfflineContent: form.OfflineContent,
	}
	exists := false
	if form.Avatar == nil {
		updateData.Avatar = 0
	} else {
		exists, err = service.File().Exists(ctx, do.CustomerChatFiles{
			Id:         form.Avatar.Id,
			CustomerId: service.AdminCtx().GetCustomerId(ctx),
			Type:       consts.FileTypeImage,
		})
		if err != nil {
			return
		}
		if !exists {
			return gerror.New("无效的图片文件")
		}
		updateData.Avatar = form.Avatar.Id
	}
	if form.Background == nil {
		updateData.Background = 0
	} else {
		exists, err = service.File().Exists(ctx, do.CustomerChatFiles{
			Id:         form.Background.Id,
			CustomerId: service.AdminCtx().GetCustomerId(ctx),
			Type:       consts.FileTypeImage,
		})
		if err != nil {
			return
		}
		if !exists {
			return gerror.New("无效的图片文件")
		}
		updateData.Background = form.Background.Id
	}
	_, err = dao.CustomerAdminChatSettings.Ctx(ctx).Data(updateData).Where("admin_id", admin.Id).Update()
	if err != nil {
		return
	}
	return
}

func (s *sAdmin) GetSetting(ctx context.Context, admin *model.CustomerAdmin) (*api.CurrentAdminSetting, error) {
	var setting *entity.CustomerAdminChatSettings
	if admin.Setting != nil {
		setting = admin.Setting
	} else {
		err := dao.CustomerAdminChatSettings.Ctx(ctx).Where("admin_id", admin.Id).Scan(&admin.Setting)
		if err != nil {
			return nil, err
		}
		if admin.Setting == nil {
			setting := &entity.CustomerAdminChatSettings{
				AdminId:        admin.Id,
				Name:           admin.Username,
				IsAutoAccept:   0,
				WelcomeContent: "",
				Avatar:         0,
			}
			result, err := dao.CustomerAdminChatSettings.Ctx(ctx).Save(*setting)
			if err != nil {
				return nil, err
			}
			id, err := result.LastInsertId()
			if err != nil {
				return nil, err
			}
			setting.Id = uint(id)
			admin.Setting = setting
			return nil, err
		}
		setting = admin.Setting
	}
	apiSetting := &api.CurrentAdminSetting{
		CurrentAdminSettingForm: api.CurrentAdminSettingForm{
			Background:     nil,
			IsAutoAccept:   setting.IsAutoAccept > 0,
			WelcomeContent: setting.WelcomeContent,
			OfflineContent: setting.OfflineContent,
			Name:           setting.Name,
			Avatar:         nil,
		},
		AdminId: admin.Id,
	}
	if setting.Background > 0 {
		apiSetting.Background, _ = service.File().FindAnd2Api(ctx, setting.Background)
	}
	if setting.Avatar > 0 {
		apiSetting.Avatar, _ = service.File().FindAnd2Api(ctx, setting.Avatar)
	}
	return apiSetting, nil

}

func (s *sAdmin) GetAvatar(ctx context.Context, model *model.CustomerAdmin) (string, error) {
	setting, err := s.GetSetting(ctx, model)
	if err != nil {
		return "", err
	}
	if setting != nil && setting.Avatar != nil {
		return setting.Avatar.Url, nil
	}
	return "", nil
}

func (s *sAdmin) GetChatName(ctx context.Context, model *model.CustomerAdmin) (string, error) {
	setting, err := s.GetSetting(ctx, model)
	if err != nil {
		return "", nil
	}
	if setting != nil && setting.Name != "" {
		return setting.Name, nil
	}
	return model.Username, nil
}
