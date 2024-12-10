package chattransfer

import (
	"context"
	"fmt"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/consts"
	"gf-chat/internal/dao"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"
	"gf-chat/internal/trait"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

func init() {
	service.RegisterChatTransfer(&sChatTransfer{
		Curd: trait.Curd[model.CustomerChatTransfer]{
			Dao: &dao.CustomerChatTransfers,
		},
	})
}

const (
	// 转接待接入的用户 sets
	transferUserKey = "user:%d:transfer"
)

type sChatTransfer struct {
	trait.Curd[model.CustomerChatTransfer]
}

func (s *sChatTransfer) ToApi(relation *model.CustomerChatTransfer) api.ChatTransfer {
	formName := ""
	toName := ""
	username := ""
	if relation.FormAdmin != nil {
		formName = relation.FormAdmin.Username
	}
	if relation.ToAdmin != nil {
		toName = relation.ToAdmin.Username
	}
	if relation.User != nil {
		username = relation.User.Username
	}
	return api.ChatTransfer{
		Id:            relation.Id,
		FromSessionId: relation.FromSessionId,
		ToSessionId:   relation.ToSessionId,
		UserId:        relation.UserId,
		Remark:        relation.Remark,
		FromAdminName: formName,
		ToAdminName:   toName,
		Username:      username,
		CreatedAt:     relation.CreatedAt,
		AcceptedAt:    relation.AcceptedAt,
		CanceledAt:    relation.CanceledAt,
	}
}

// Cancel 取消待接入的转接
func (s *sChatTransfer) Cancel(ctx context.Context, transfer *model.CustomerChatTransfer) (err error) {
	if transfer.CanceledAt != nil {
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, "该转接已取消")
	}
	if transfer.AcceptedAt != nil {
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, "该转接已接入，无法取消")
	}
	now := gtime.Now()
	if transfer.ToSession != nil {
		_, err = service.ChatSession().UpdatePri(ctx, transfer.ToSessionId, do.CustomerChatSessions{
			CanceledAt: now,
		})
		if err != nil {
			return
		}
	}
	_, err = s.UpdatePri(ctx, transfer.Id, do.CustomerChatTransfers{
		CanceledAt: now,
	})
	if err != nil {
		return
	}
	err = s.removeUser(ctx, transfer.CustomerId, transfer.UserId)
	if err != nil {
		return
	}
	err = service.Chat().NoticeTransfer(ctx, transfer.CustomerId, transfer.ToAdminId)
	return
}

// Accept 接入转接
func (s *sChatTransfer) Accept(ctx context.Context, transfer *model.CustomerChatTransfer) error {
	if transfer.AcceptedAt != nil {
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, "该转接已被接入")
	}
	if transfer.CanceledAt != nil {
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, "该转接已取消")
	}
	_, err := s.UpdatePri(ctx, transfer.Id, do.CustomerChatTransfers{
		AcceptedAt: gtime.Now(),
	})
	if err != nil {
		return err
	}
	err = service.Chat().NoticeTransfer(ctx, transfer.CustomerId, transfer.ToAdminId)
	if err != nil {
		return err
	}
	return s.removeUser(ctx, transfer.CustomerId, transfer.UserId)
}

// Create 创建转接
func (s *sChatTransfer) Create(ctx context.Context, fromAdminId, toId, uid uint, remark string) (err error) {
	session, err := service.ChatSession().FirstActive(ctx, uid, fromAdminId, nil)
	if err != nil {
		return err
	}
	if session == nil {
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, "用户已失效，无法转接")
	}
	err = service.ChatSession().Close(ctx, session, true, false)
	if err != nil {
		return
	}
	newSession := &model.CustomerChatSession{
		CustomerChatSessions: entity.CustomerChatSessions{
			QueriedAt:  gtime.Now(),
			CustomerId: session.CustomerId,
			AdminId:    toId,
			Type:       consts.ChatSessionTypeTransfer,
			UserId:     uid,
		},
	}
	newSession, err = service.ChatSession().Insert(ctx, newSession)
	if err != nil {
		return err
	}
	transfer := &model.CustomerChatTransfer{
		CustomerChatTransfers: entity.CustomerChatTransfers{
			UserId:        uid,
			FromSessionId: session.Id,
			ToSessionId:   newSession.Id,
			CustomerId:    session.CustomerId,
			FromAdminId:   fromAdminId,
			ToAdminId:     toId,
			Remark:        remark,
		},
	}
	_, err = service.ChatTransfer().Save(ctx, transfer)
	if err != nil {
		return
	}
	err = s.addUser(ctx, session.CustomerId, uid, toId)
	if err != nil {
		return
	}
	err = service.Chat().NoticeTransfer(ctx, session.CustomerId, toId)
	if err != nil {
		return
	}
	return nil
}
func (s *sChatTransfer) userKey(customerId uint) string {
	return fmt.Sprintf(transferUserKey, customerId)
}

// RemoveUser 在转接列表中移除user
func (s *sChatTransfer) removeUser(ctx context.Context, customerId, uid uint) error {
	_, err := g.Redis().Do(ctx, "hDel", s.userKey(customerId), uid)
	return err
}

func (s *sChatTransfer) GetUserTransferId(ctx context.Context, customerId, uid uint) (id uint, error error) {
	val, err := g.Redis().Do(ctx, "hGet", s.userKey(customerId), uid)
	if err != nil {
		return
	}
	return val.Uint(), nil
}

// IsInTransfer 是否待转接
func (s *sChatTransfer) IsInTransfer(ctx context.Context, customerId uint, uid uint) (is bool, err error) {
	id, err := s.GetUserTransferId(ctx, customerId, uid)
	if err != nil {
		return
	}
	return id != 0, nil
}

// AddUser 添加用户到转接列表中
func (s *sChatTransfer) addUser(ctx context.Context, customerId, uid, adminId uint) (err error) {
	_, err = g.Redis().Do(ctx, "hSet", s.userKey(customerId), uid, adminId)
	return
}
