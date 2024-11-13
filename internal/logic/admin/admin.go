package admin

import (
	"context"
	"gf-chat/internal/dao"
	"gf-chat/internal/model"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"
	"gf-chat/internal/trait"
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

func (s *sAdmin) CanAccess(admin *model.CustomerAdmin) bool {
	return true
}

func (s *sAdmin) GetSetting(ctx context.Context, admin *model.CustomerAdmin) (*entity.CustomerAdminChatSettings, error) {
	if admin.Setting != nil {
		return admin.Setting, nil
	}
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
			Avatar:         "",
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
	return admin.Setting, nil
}

func (s *sAdmin) GetAvatar(ctx context.Context, model *model.CustomerAdmin) (string, error) {
	setting, err := s.GetSetting(ctx, model)
	if err != nil {
		return "", err
	}
	if setting.Avatar != "" {
		//return service.Qiniu().Url(model.Setting.Avatar), nil
	} else {
		return "", nil
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
