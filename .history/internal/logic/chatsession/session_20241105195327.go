package chat

import (
	"context"
	"database/sql"
	"gf-chat/internal/consts"
	"gf-chat/internal/dao"
	"gf-chat/internal/model"
	"gf-chat/internal/model/chat"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/model/relation"
	"gf-chat/internal/service"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

func init() {
	service.RegisterChatSession(&sChatSession{})
}

type sChatSession struct {
}

func (s *sChatSession) Get(ctx context.Context, w any) (res []*relation.CustomerChatSessions) {
	err := dao.CustomerChatSessions.Ctx(ctx).Where(w).WithAll().OrderDesc("id").Scan(&res)
	if err == sql.ErrNoRows {
		return
	}
	return
}

func (s *sChatSession) Paginate(ctx context.Context, w any, page model.QueryInput) (res []*relation.CustomerChatSessions, total int) {
	query := dao.CustomerChatSessions.Ctx(ctx).Where(w)
	if page.WithTotal {
		total, _ = query.Clone().Count()
		if total == 0 {
			return
		}
	}
	err := query.Page(page.Page, page.Size).WithAll().OrderDesc("id").Scan(&res)
	if err == sql.ErrNoRows {
		return
	}
	return
}

func (s *sChatSession) Cancel(session *entity.CustomerChatSessions) error {
	if session.AcceptedAt != nil {
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, "会话已接入，无法取消")
	}
	if session.CanceledAt != nil {
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, "会话已取消，请勿重复取消")
	}
	session.CanceledAt = gtime.New()
	dao.CustomerChatSessions.Ctx(gctx.New()).Save(session)
	service.ChatManual().Remove(session.UserId, session.CustomerId)
	service.Chat().BroadcastWaitingUser(session.CustomerId)
	return nil
}

// Close 关闭会话
func (s *sChatSession) Close(session *entity.CustomerChatSessions, isRemoveUser bool, updateTime bool) {
	if session.BrokenAt != nil {
		session.BrokenAt = gtime.New()
		dao.CustomerChatSessions.Ctx(gctx.New()).Save(session)
	}
	if isRemoveUser {
		service.ChatRelation().RemoveUser(session.AdminId, session.UserId)
	} else {
		if updateTime {
			service.ChatRelation().UpdateLimitTime(session.AdminId, session.UserId, 0)
		}
	}
}

func (s *sChatSession) RelationToChat(model *relation.CustomerChatSessions) chat.Session {
	adminName := ""
	if model.Admin != nil {
		adminName = model.Admin.Username
	}
	username := ""
	if model.User != nil {
		username = model.User.Username
	}
	statusLabel := ""
	status := ""
	if model.CanceledAt != nil {
		statusLabel = "已取消"
		status = consts.ChatSessionStatusCancel
	}
	if model.CanceledAt == nil && model.AcceptedAt == nil {
		statusLabel = "待接入"
		status = consts.ChatSessionStatusWait
	}
	if model.AcceptedAt != nil && model.BrokenAt == nil {
		statusLabel = "已接入"
		status = consts.ChatSessionStatusAccept
	}
	if model.BrokenAt != nil {
		statusLabel = "已关闭"
		status = consts.ChatSessionStatusClose
	}
	typeLabel := ""
	if model.Type == consts.ChatSessionTypeNormal {
		typeLabel = "普通"
	}
	if model.Type == consts.ChatSessionTypeTransfer {
		typeLabel = "转接"
	}

	return chat.Session{
		Id:          model.Id,
		UserId:      model.UserId,
		QueriedAt:   model.QueriedAt,
		AcceptedAt:  model.AcceptedAt,
		BrokeAt:     model.BrokeAt,
		CanceledAt:  model.CanceledAt,
		AdminId:     model.AdminId,
		UserName:    username,
		AdminName:   adminName,
		TypeLabel:   typeLabel,
		Status:      status,
		StatusLabel: statusLabel,
		Rate:        model.Rate,
	}
}
func (s *sChatSession) FirstRelation(ctx context.Context, w do.CustomerChatSessions) *relation.CustomerChatSessions {
	session := &relation.CustomerChatSessions{}
	err := dao.CustomerChatSessions.Ctx(ctx).Where(w).WithAll().Scan(session)
	if err == sql.ErrNoRows {
		return nil
	}
	return session
}
func (s *sChatSession) First(ctx context.Context, w do.CustomerChatSessions) *entity.CustomerChatSessions {
	session := &entity.CustomerChatSessions{}
	err := dao.CustomerChatSessions.Ctx(ctx).Where(w).Scan(session)
	if err == sql.ErrNoRows {
		return nil
	}
	return session
}

func (s *sChatSession) SaveEntity(model *entity.CustomerChatSessions) *entity.CustomerChatSessions {
	result, _ := dao.CustomerChatSessions.Ctx(gctx.New()).Save(model)
	id, _ := result.LastInsertId()
	model.Id = gconv.Uint64(id)
	return model
}

func (s *sChatSession) Create(uid uint, customerId uint, t uint) *entity.CustomerChatSessions {
	model := &entity.CustomerChatSessions{
		UserId:     uid,
		Type:       t,
		CustomerId: customerId,
		QueriedAt:  time.Now().Unix(),
	}
	return s.SaveEntity(model)
}

func (s *sChatSession) GetUnAcceptModel(customerId uint) (res []*relation.CustomerChatSessions) {
	dao.CustomerChatSessions.Ctx(gctx.New()).Where(do.CustomerChatSessions{
		CanceledAt: 0,
		AdminId:    0,
		Type:       consts.ChatSessionTypeNormal,
		CustomerId: customerId,
	}).WithAll().Scan(&res)
	return
}
func (s *sChatSession) ActiveTransferOne(uid uint, adminId uint) *entity.CustomerChatSessions {
	return s.ActiveOne(uid, adminId, consts.ChatSessionTypeTransfer)
}

func (s *sChatSession) ActiveNormalOne(uid uint, adminId uint) *entity.CustomerChatSessions {
	return s.ActiveOne(uid, adminId, consts.ChatSessionTypeNormal)
}

func (s *sChatSession) ActiveOne(uid uint, adminId, t any) *entity.CustomerChatSessions {
	w := do.CustomerChatSessions{
		UserId:     uid,
		AdminId:    adminId,
		CanceledAt: 0,
		BrokeAt:    0,
	}
	if t != nil {
		w.Type = t
	}
	return s.First(gctx.New(), w)
}
