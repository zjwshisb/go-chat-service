package chat

import (
	"context"
	api "gf-chat/api/backend/v1"
	"gf-chat/internal/consts"
	"gf-chat/internal/dao"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"
	"gf-chat/internal/trait"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
)

func init() {
	service.RegisterChatSession(newSChatSession())
}

func newSChatSession() *sChatSession {
	return &sChatSession{
		trait.Curd[model.CustomerChatSession]{
			Dao: &dao.CustomerChatSessions,
		},
	}
}

type sChatSession struct {
	trait.Curd[model.CustomerChatSession]
}

func (s *sChatSession) Cancel(ctx context.Context, session *model.CustomerChatSession) (err error) {
	if session.AcceptedAt != nil {
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, "会话已接入，无法取消")
	}
	if session.CanceledAt != nil {
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, "会话已取消，请勿重复取消")
	}
	_, err = s.UpdatePri(ctx, session.Id, do.CustomerChatSessions{
		CanceledAt: gtime.Now(),
	})
	if err != nil {
		return
	}
	err = service.Chat().RemoveManual(ctx, session.UserId, session.CustomerId)
	if err != nil {
		return
	}
	_ = service.Chat().BroadcastWaitingUser(ctx, session.CustomerId)
	return nil
}

// Close 关闭会话
func (s *sChatSession) Close(ctx context.Context, session *model.CustomerChatSession, isRemoveUser bool, updateTime bool) (err error) {
	if session.AcceptedAt == nil {
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, "未接入会话无法断开")
	}
	if session.BrokenAt != nil {
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, "会话关闭，请勿重复操作")
	}
	_, err = s.UpdatePri(ctx, session.Id, do.CustomerChatSessions{
		BrokenAt: gtime.Now(),
	})
	if err != nil {
		return
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

func (s *sChatSession) ToApi(session *model.CustomerChatSession) api.ChatSession {
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
func (s *sChatSession) Create(ctx context.Context, uid uint, customerId uint, t uint) (item *model.CustomerChatSession, err error) {
	item = &model.CustomerChatSession{}
	item.Type = t
	item.CustomerId = customerId
	item.QueriedAt = gtime.Now()
	item.UserId = uid
	id, err := s.Save(ctx, item)
	if err != nil {
		return
	}
	item.CustomerChatSessions.Id = uint(id)
	return item, nil
}

func (s *sChatSession) GetUnAccepts(ctx context.Context, customerId uint) (res []*model.CustomerChatSession, err error) {
	return s.All(ctx, g.Map{
		"canceled_at is null": nil,
		"admin_id":            0,
		"type":                consts.ChatSessionTypeNormal,
		"customer_id":         customerId,
	}, g.Slice{model.CustomerChatSession{}.User}, nil)

}
func (s *sChatSession) Insert(ctx context.Context, session *model.CustomerChatSession) (m *model.CustomerChatSession, err error) {
	id, err := s.Save(ctx, session)
	if err != nil {
		return
	}
	session.Id = gconv.Uint(id)
	return session, nil
}

func (s *sChatSession) FirstTransfer(ctx context.Context, uid uint, adminId uint) (*model.CustomerChatSession, error) {
	return s.FirstActive(ctx, uid, adminId, consts.ChatSessionTypeTransfer)
}

func (s *sChatSession) FirstNormal(ctx context.Context, uid uint, adminId uint) (*model.CustomerChatSession, error) {
	return s.FirstActive(ctx, uid, adminId, consts.ChatSessionTypeNormal)
}

func (s *sChatSession) FirstActive(ctx context.Context, uid uint, adminId, t any) (*model.CustomerChatSession, error) {

	w := g.Map{
		"user_id":             uid,
		"admin_id":            adminId,
		"canceled_at is null": nil,
		"broken_at is null":   nil,
	}
	if t != nil {
		w["type"] = t
	}
	return s.First(ctx, w)
}
