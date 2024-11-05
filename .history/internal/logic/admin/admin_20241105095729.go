package admin

import (
	"context"
	"database/sql"
	"errors"
	"gf-chat/internal/dao"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/model/relation"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/golang-module/carbon/v2"
	"strconv"
)

func init() {
	service.RegisterAdmin(&sAdmin{})
}

type sAdmin struct {
}

func (s *sAdmin) GetAdmins(ctx g.Ctx, w any) (items []*entity.CustomerAdmins) {
	err := dao.CustomerAdmins.Ctx(ctx).Where(w).Scan(&items)
	if err == sql.ErrNoRows {
		return
	}
	return
}
func (s *sAdmin) Paginate(ctx context.Context, where *do.CustomerAdmins, p model.QueryInput) (items []*relation.CustomerAdmins, total int) {
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

func (s *sAdmin) IsValid(admin *entity.CustomerAdmins) error {
	if admin.Id == 0 {
		return errors.New("没有权限登录")
	}
	if admin.Status == 0 {
		return errors.New("没有权限登录")
	}
	if admin.IsChat == 0 {
		return errors.New("没有客服权限")
	}
	setting := &entity.CustomerSettings{}
	_ = dao.CustomerSettings.Ctx(gctx.New()).Where("customer_id", admin.CustomerId).Scan(setting)
	if setting.ChatOpen == 0 {
		return errors.New("没有该系统的权限")
	}
	return nil
}

func (s *sAdmin) EntityToRelation(admin *entity.CustomerAdmins) *relation.CustomerAdmins {
	setting := s.GetSetting(admin.Id)
	return &relation.CustomerAdmins{
		CustomerAdmins: *admin,
		Setting:        setting,
	}
}

func (s *sAdmin) GetSetting(adminId uint) *entity.CustomerAdminChatSettings {
	setting := &entity.CustomerAdminChatSettings{}
	ctx := gctx.New()
	err := dao.CustomerAdminChatSettings.Ctx(ctx).Where("admin_id", adminId).Scan(setting)
	if err == sql.ErrNoRows {
		setting.AdminId = adminId
		result, _ := dao.CustomerAdminChatSettings.Ctx(ctx).Save(setting)
		id, _ := result.LastInsertId()
		setting.Id = id
		return setting
	}
	return setting
}

func (s *sAdmin) GetAvatar(model *relation.CustomerAdmins) string {
	if model.Setting != nil && model.Setting.Avatar != "" {
		return service.Qiniu().Url(model.Setting.Avatar)
	} else {
		return ""
	}
}

func (s *sAdmin) GetChatName(model *entity.CustomerAdmins) string {
	setting := s.GetSetting(model.Id)
	if setting != nil && setting.Name != "" {
		return setting.Name
	}
	return model.Username
}

func (s *sAdmin) First(id int) (admin *entity.CustomerAdmins) {
	admin = &entity.CustomerAdmins{}
	err := dao.CustomerAdmins.Ctx(gctx.New()).WherePri(id).Scan(admin)
	if err != nil {
		return nil
	}
	return
}
func (s *sAdmin) FirstRelation(id int) *relation.CustomerAdmins {
	admin := s.First(id)
	if admin != nil {
		setting := s.GetSetting(admin.Id)
		return &relation.CustomerAdmins{
			CustomerAdmins: *admin,
			Setting:        setting,
		}
	}
	return nil
}

func (s *sAdmin) GetWechat(adminId uint) *entity.CustomerAdminWechat {
	wechat := &entity.CustomerAdminWechat{}
	err := dao.CustomerAdminWechat.Ctx(gctx.New()).Where("admin_id", adminId).Scan(wechat)
	if err == sql.ErrNoRows {
		return nil
	}
	return wechat
}

func (s *sAdmin) GetChatAll(customerId int) []*relation.CustomerAdmins {
	admins := make([]*relation.CustomerAdmins, 0)
	_ = dao.CustomerAdmins.Ctx(gctx.New()).Where(do.CustomerAdmins{
		CustomerId: customerId,
		IsChat:     1,
		Status:     1,
	}).WithAll().Scan(&admins)
	return admins
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
	sessions := make([]*entity.CustomerChatSessions, 0, 0)

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
			d := (session.AcceptedAt - firstDateUnix) / (24 * 3600)
			lines[d].Value += 1
		}
	}

	return lines, &model.AdminDetailInfo{
		Avatar:        "",
		Username:      admin.Username,
		Online:        service.Chat().IsOnline(admin.CustomerId, gconv.Int(admin.Id), "admin"),
		Id:            admin.Id,
		AcceptedCount: service.ChatRelation().GetActiveCount(gconv.Int(admin.Id)),
	}, nil
}
