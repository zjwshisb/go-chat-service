package admin

import (
	"context"
	"database/sql"
	"errors"
	api "gf-chat/api/backend/v1"
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
	"github.com/gogf/gf/v2/os/gtime"
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

// Login authenticates an admin user with username and password, returning the admin model and JWT token if successful.
// Returns error if authentication fails or admin account is disabled.
func (s *sAdmin) Login(ctx context.Context, request *ghttp.Request) (admin *model.CustomerAdmin, token string, err error) {
	username := request.Get("username")
	password := request.Get("password")
	admin, err = s.First(ctx, do.CustomerAdmins{Username: username.String()})
	if err != nil {
		err = gerror.NewCode(gcode.CodeBusinessValidationFailed, "账号或密码错误")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), password.Bytes())
	if err != nil {
		err = gerror.NewCode(gcode.CodeBusinessValidationFailed, "账号或密码错误")
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

// Auth authenticates a user based on their JWT token, returning the admin model and error if authentication fails.
// Returns error if token is missing or invalid.
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

// FindSetting retrieves a customer admin chat setting by admin ID, optionally with file details.
// Returns error if no setting is found or database query fails.
func (s *sAdmin) FindSetting(ctx context.Context, adminId uint, withFile bool) (*model.CustomerAdminChatSetting, error) {
	var setting *model.CustomerAdminChatSetting
	query := dao.CustomerAdminChatSettings.Ctx(ctx).Where("admin_id", adminId)
	if withFile {
		query = query.WithAll()
	}
	err := query.Scan(&setting)
	if err != nil {
		return nil, err
	}
	if setting == nil {
		return nil, sql.ErrNoRows
	}
	return setting, nil
}

// GenSetting generates a new customer admin chat setting for a given admin.
// Returns error if database save fails.
func (s *sAdmin) GenSetting(ctx context.Context, admin *model.CustomerAdmin) (*model.CustomerAdminChatSetting, error) {
	setting := &model.CustomerAdminChatSetting{
		CustomerAdminChatSettings: entity.CustomerAdminChatSettings{
			AdminId:        admin.Id,
			Name:           admin.Username,
			IsAutoAccept:   0,
			WelcomeContent: "",
			Avatar:         0,
		},
	}
	result, err := dao.CustomerAdminChatSettings.Ctx(ctx).Save(*setting)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	setting.Id = gconv.Uint(id)
	return setting, nil
}

// UpdateLastOnline updates the last online timestamp for a customer admin.
// Returns error if database update fails.
func (s *sAdmin) UpdateLastOnline(ctx context.Context, adminId uint) error {
	_, err := dao.CustomerAdminChatSettings.Ctx(ctx).Where("admin_id", adminId).Update(do.CustomerAdminChatSettings{
		LastOnline: gtime.Now(),
	})
	return err
}

// UpdateSetting updates the chat setting for a customer admin.
// Returns error if database update fails.
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
			return gerror.NewCode(gcode.CodeValidationFailed, "无效的图片文件")
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
			return gerror.NewCode(gcode.CodeValidationFailed, "无效的图片文件")
		}
		updateData.Background = form.Background.Id
	}
	_, err = dao.CustomerAdminChatSettings.Ctx(ctx).Data(updateData).Where("admin_id", admin.Id).Update()
	if err != nil {
		return
	}
	return
}

// GetAndFindSetting retrieves a customer admin chat setting, either from the admin's setting or from the database.
// Returns error if database query fails.
func (s *sAdmin) GetAndFindSetting(ctx context.Context, admin *model.CustomerAdmin) (*model.CustomerAdminChatSetting, error) {
	if admin.Setting != nil {
		return admin.Setting, nil
	} else {
		return s.FindSetting(ctx, admin.Id, true)
	}
}

// GetApiSetting retrieves a customer admin chat setting and converts it to an API format.
// Returns error if database query fails.
func (s *sAdmin) GetApiSetting(ctx context.Context, admin *model.CustomerAdmin) (*api.CurrentAdminSetting, error) {
	var setting *model.CustomerAdminChatSetting
	if admin.Setting != nil {
		setting = admin.Setting
	} else {
		err := dao.CustomerAdminChatSettings.Ctx(ctx).Where("admin_id", admin.Id).Scan(&admin.Setting)
		if err != nil {
			return nil, err
		}
		if admin.Setting == nil {
			setting, err := s.GenSetting(ctx, admin)
			if err != nil {
				return nil, err
			}
			admin.Setting = setting
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
	if setting.BackgroundFile != nil {
		apiSetting.Background = service.File().ToApi(setting.BackgroundFile)
	} else if setting.Background > 0 {
		apiSetting.Background, _ = service.File().FindAnd2Api(ctx, setting.Background)
	}
	if setting.AvatarFile != nil {
		apiSetting.Avatar = service.File().ToApi(setting.AvatarFile)
	} else if setting.Avatar > 0 {
		apiSetting.Avatar, _ = service.File().FindAnd2Api(ctx, setting.Avatar)
	}
	return apiSetting, nil
}

// GetAvatar retrieves the avatar URL for a customer admin, either from the admin's setting or from the database.
// Returns error if retrieval fails.
func (s *sAdmin) GetAvatar(ctx context.Context, admin *model.CustomerAdmin) (string, error) {
	setting, err := s.GetAndFindSetting(ctx, admin)
	if err != nil {
		return "", err
	}
	if setting != nil {
		if setting.AvatarFile != nil {
			return service.File().Url(setting.AvatarFile), nil
		} else if setting.Avatar > 0 {
			avatarFile, err := service.File().Find(ctx, setting.Avatar)
			if err != nil {
				return "", err
			}
			return service.File().Url(avatarFile), nil
		}
	}
	return "", nil
}

// GetChatName retrieves the chat name for a customer admin, either from the admin's setting or from the database.
// Returns error if retrieval fails.
func (s *sAdmin) GetChatName(ctx context.Context, admin *model.CustomerAdmin) (string, error) {
	setting, err := s.GetAndFindSetting(ctx, admin)
	if err != nil {
		return "", nil
	}
	if setting != nil && setting.Name != "" {
		return setting.Name, nil
	}
	return admin.Username, nil
}

// GetAdminsWithSetting retrieves all customer admins with their settings, optionally filtered by a where clause.
// Returns error if database query fails.
func (s *sAdmin) GetAdminsWithSetting(ctx context.Context, where any) (res []*model.CustomerAdmin, err error) {
	err = s.Dao.Ctx(ctx).Where(where).WithAll().Scan(&res)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return
	}
	if res == nil {
		res = make([]*model.CustomerAdmin, 0)
	}
	return
}
