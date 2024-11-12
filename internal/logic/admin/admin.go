package admin

import (
	"context"
	"database/sql"
	"gf-chat/internal/dao"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"
	"gf-chat/internal/trait"
	"strconv"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/golang-module/carbon/v2"
)

func init() {
	service.RegisterAdmin(&sAdmin{
		trait.Curd[model.CustomerAdmin]{
			Dao: dao.CustomerAdmins,
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
		return service.Qiniu().Url(model.Setting.Avatar), nil
	} else {
		return "", nil
	}
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

func (s *sAdmin) GetDetail(ctx context.Context, id any, month string) ([]*model.ChartLine, *model.AdminDetailInfo, error) {
	admin, err := s.First(ctx, do.CustomerAdmins{
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
	})
	if err != nil {
		return nil, nil, err
	}

	date := carbon.Parse(month)
	if date.Error != nil {
		date = carbon.Now()
	}
	firstDate := date.StartOfMonth()
	lastDate := date.EndOfMonth()
	sessions := make([]*entity.CustomerChatSessions, 0)

	err = dao.CustomerChatSessions.Ctx(ctx).Where(g.Map{
		"accepted_at>=": firstDate.ToDateTimeString(),
		"accepted_at<=": lastDate.ToDateTimeString(),
		"admin_id":      admin.Id,
	}).Scan(&sessions)
	lines := make([]*model.ChartLine, lastDate.DayOfMonth())
	if err != sql.ErrNoRows {
		for day := range lines {
			lines[day] = &model.ChartLine{
				Category: "每日接待数",
				Value:    0,
				Label:    strconv.Itoa(day+1) + "号",
			}
		}
		for _, session := range sessions {
			d := (session.AcceptedAt.Unix() - firstDate.Carbon2Time().Unix()) / (24 * 3600)
			lines[d].Value += 1
		}
	}

	return lines, &model.AdminDetailInfo{
		Avatar:        "",
		Username:      admin.Username,
		Online:        service.Chat().IsOnline(admin.CustomerId, admin.Id, "admin"),
		Id:            admin.Id,
		AcceptedCount: service.ChatRelation().GetActiveCount(ctx, admin.Id),
	}, nil
}
