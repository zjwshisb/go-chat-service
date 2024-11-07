package admin

import (
	"context"
	"database/sql"
	"errors"
	"gf-chat/internal/dao"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"
	"strconv"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/golang-module/carbon/v2"
)

func init() {
	service.RegisterAdmin(&sAdmin{})
}

type sAdmin struct {
}

func (s *sAdmin) All(ctx g.Ctx, w do.CustomerAdmins, with ...any) (items []*model.CustomerAdmin, err error) {
	err = dao.CustomerAdmins.Ctx(ctx).Where(w).Scan(&items)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = make([]*model.CustomerAdmin, 0)
	}
	return
}
func (s *sAdmin) Paginate(ctx context.Context, where *do.CustomerAdmins, p model.QueryInput) (items []*model.CustomerAdmin, total int) {
	query := dao.CustomerAdmins.Ctx(ctx)
	if where != nil {
		query = query.Where(where)
	}
	if p.WithTotal {
		total, _ = query.Clone().Count()
		if total == 0 {
			return
		}
	}
	query = query.Page(p.GetPage(), p.GetSize())
	err := query.WithAll().Unscoped().Scan(&items)
	if err == sql.ErrNoRows {
		return
	}
	return
}

func (s *sAdmin) IsValid(admin *model.CustomerAdmin) error {
	if admin == nil {
		return errors.New("没有权限登录")
	}

	return nil
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
		err = gerror.New("no setting")
	}
	return admin.Setting, nil
}

func (s *sAdmin) GetAvatar(model *model.CustomerAdmin) (string, error) {
	if model.Setting != nil && model.Setting.Avatar != "" {
		return service.Qiniu().Url(model.Setting.Avatar), nil
	} else {
		return "", gerror.New("no avatar")
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

func (s *sAdmin) First(ctx context.Context, where do.CustomerAdmins) (admin *model.CustomerAdmin, err error) {
	err = dao.CustomerAdmins.Ctx(ctx).Where(where).Scan(&admin)
	if err != nil {
		return
	}
	if admin == nil {
		err = sql.ErrNoRows
	}
	return
}

func (s *sAdmin) GetWechat(adminId uint) *entity.CustomerAdminWechat {
	wechat := &entity.CustomerAdminWechat{}
	err := dao.CustomerAdminWechat.Ctx(gctx.New()).Where("admin_id", adminId).Scan(wechat)
	if err == sql.ErrNoRows {
		return nil
	}
	return wechat
}

func (s *sAdmin) GetDetail(ctx context.Context, id any, month string) ([]*model.ChartLine, *model.AdminDetailInfo, error) {
	admin := &entity.CustomerAdmins{}
	err := dao.CustomerAdmins.Ctx(ctx).
		Where("customer_id", service.AdminCtx().GetCustomerId(ctx)).WherePri(id).Scan(admin)
	if err != nil {
		return nil, nil, err
	}
	date := carbon.Parse(month)
	if date.Error != nil {
		date = carbon.Now()
	}
	firstDate := date.StartOfMonth()
	lastDate := date.EndOfMonth()
	firstDateUnix := firstDate.Timestamp()
	sessions := make([]*entity.CustomerChatSessions, 0)

	err = dao.CustomerChatSessions.Ctx(ctx).Where(g.Map{
		"accepted_at>=": firstDateUnix,
		"accepted_at<=": lastDate.Timestamp(),
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
			d := (session.AcceptedAt.Unix() - firstDateUnix) / (24 * 3600)
			lines[d].Value += 1
		}
	}

	return lines, &model.AdminDetailInfo{
		Avatar:        "",
		Username:      admin.Username,
		Online:        service.Chat().IsOnline(admin.CustomerId, admin.Id, "admin"),
		Id:            admin.Id,
		AcceptedCount: service.ChatRelation().GetActiveCount(admin.Id),
	}, nil
}
