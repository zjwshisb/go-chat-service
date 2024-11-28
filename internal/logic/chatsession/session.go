package chat

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
	"github.com/gogf/gf/v2/frame/g"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
)

func init() {
	service.RegisterChatSession(new())
}

func new() *sChatSession {
	return &sChatSession{
		trait.Curd[model.CustomerChatSession]{
			Dao: &dao.CustomerChatSessions,
		},
	}
}

type sChatSession struct {
	trait.Curd[model.CustomerChatSession]
}

// func (s *sChatSession) Get(ctx context.Context, w any) (res []*model.CustomerChatSession) {
// 	err := dao.CustomerChatSessions.Ctx(ctx).Where(w).WithAll().OrderDesc("id").Scan(&res)
// 	if err == sql.ErrNoRows {
// 		return
// 	}
// 	return
// }

// func (s *sChatSession) Paginate(ctx context.Context, w any, page model.QueryInput) (res []*model.CustomerChatSession, total int) {
// 	query := dao.CustomerChatSessions.Ctx(ctx).Where(w)
// 	if page.WithTotal {
// 		total, _ = query.Clone().Count()
// 		if total == 0 {
// 			return
// 		}
// 	}
// 	err := query.Page(page.Page, page.Size).WithAll().OrderDesc("id").Scan(&res)
// 	if err == sql.ErrNoRows {
// 		return
// 	}
// 	return
// }

func (s *sChatSession) Cancel(ctx context.Context, session *model.CustomerChatSession) (err error) {
	if session.AcceptedAt != nil {
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, "会话已接入，无法取消")
	}
	if session.CanceledAt != nil {
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, "会话已取消，请勿重复取消")
	}
	session.CanceledAt = gtime.New()
	_, err = s.Save(ctx, session)
	if err != nil {
		return
	}
	err = service.Chat().RemoveManual(ctx, session.UserId, session.CustomerId)
	if err != nil {
		return
	}
	service.Chat().BroadcastWaitingUser(ctx, session.CustomerId)
	return nil
}

// Close 关闭会话
func (s *sChatSession) Close(ctx context.Context, session *model.CustomerChatSession, isRemoveUser bool, updateTime bool) (err error) {
	if session.BrokenAt != nil {
		session.BrokenAt = gtime.New()
		_, err = s.Save(ctx, session)
		if err != nil {
			return
		}
	}
	if isRemoveUser {
		err = service.ChatRelation().RemoveUser(ctx, session.AdminId, session.UserId)
	} else {
		if updateTime {
			err = service.ChatRelation().UpdateLimitTime(ctx, session.AdminId, session.UserId, 0)
		}
	}
	return
}

func (s *sChatSession) RelationToChat(session *model.CustomerChatSession) api.ChatSession {
	adminName := ""
	if session.Admin != nil {
		adminName = session.Admin.Username
	}
	username := ""
	if session.User != nil {
		username = session.User.Username
	}
	statusLabel := ""
	status := ""
	if session.CanceledAt != nil {
		statusLabel = "已取消"
		status = consts.ChatSessionStatusCancel
	}
	if session.CanceledAt == nil && session.AcceptedAt == nil {
		statusLabel = "待接入"
		status = consts.ChatSessionStatusWait
	}
	if session.AcceptedAt != nil && session.BrokenAt == nil {
		statusLabel = "已接入"
		status = consts.ChatSessionStatusAccept
	}
	if session.BrokenAt != nil {
		statusLabel = "已关闭"
		status = consts.ChatSessionStatusClose
	}
	typeLabel := ""
	if session.Type == consts.ChatSessionTypeNormal {
		typeLabel = "普通"
	}
	if session.Type == consts.ChatSessionTypeTransfer {
		typeLabel = "转接"
	}

	return api.ChatSession{
		Id:          session.Id,
		UserId:      session.UserId,
		QueriedAt:   session.QueriedAt,
		AcceptedAt:  session.AcceptedAt,
		BrokenAt:    session.BrokenAt,
		CanceledAt:  session.CanceledAt,
		AdminId:     session.AdminId,
		UserName:    username,
		AdminName:   adminName,
		TypeLabel:   typeLabel,
		Status:      status,
		StatusLabel: statusLabel,
		Rate:        session.Rate,
	}
}

// func (s *sChatSession) FirstRelation(ctx context.Context, w do.CustomerChatSessions) *model.CustomerChatSession {
// 	session := &model.CustomerChatSession{}
// 	err := dao.CustomerChatSessions.Ctx(ctx).Where(w).WithAll().Scan(session)
// 	if err == sql.ErrNoRows {
// 		return nil
// 	}
// 	return session
// }
// func (s *sChatSession) First(ctx context.Context, w do.CustomerChatSessions) (item *model.CustomerChatSession, err error) {
// 	err = dao.CustomerChatSessions.Ctx(ctx).Where(w).Scan(&item)
// 	if err != sql.ErrNoRows {
// 		return
// 	}
// 	if item == nil {
// 		err = sql.ErrNoRows
// 	}
// 	return
// }

// func (s *sChatSession) SaveEntity(ctx context.Context, model *entity.CustomerChatSessions) *entity.CustomerChatSessions {
// 	result, _ := dao.CustomerChatSessions.Ctx(ctx).Save(model)
// 	id, _ := result.LastInsertId()
// 	model.Id = uint(id)
// 	return model
// }

func (s *sChatSession) Create(ctx context.Context, uid uint, customerId uint, t uint) (item *model.CustomerChatSession, err error) {
	item = &model.CustomerChatSession{
		CustomerChatSessions: entity.CustomerChatSessions{
			UserId:     uid,
			Type:       t,
			CustomerId: customerId,
			QueriedAt:  gtime.New(),
		},
	}
	id, err := s.Save(ctx, item)
	if err != nil {
		return
	}
	item.CustomerChatSessions.Id = uint(id)
	return item, nil
}

func (s *sChatSession) GetUnAccepts(ctx context.Context, customerId uint) (res []*model.CustomerChatSession, err error) {
	return s.All(ctx, do.CustomerChatSessions{
		CanceledAt: nil,
		AdminId:    0,
		Type:       consts.ChatSessionTypeNormal,
		CustomerId: customerId,
	}, g.Slice{model.CustomerChatSession{}.User, model.CustomerChatSession{}.Admin}, nil)

}
func (s *sChatSession) FirstTransfer(ctx context.Context, uid uint, adminId uint) (*model.CustomerChatSession, error) {
	return s.FirstActive(ctx, uid, adminId, consts.ChatSessionTypeTransfer)
}

func (s *sChatSession) FirstNormal(ctx context.Context, uid uint, adminId uint) (*model.CustomerChatSession, error) {
	return s.FirstActive(ctx, uid, adminId, consts.ChatSessionTypeNormal)
}

func (s *sChatSession) FirstActive(ctx context.Context, uid uint, adminId, t any) (*model.CustomerChatSession, error) {
	w := g.Map{
		"user_id":     uid,
		"admin_id":    adminId,
		"canceled_at": nil,
		"broken_at":   nil,
	}
	if t != nil {
		w["type"] = t
	}
	return s.First(ctx, w)
}