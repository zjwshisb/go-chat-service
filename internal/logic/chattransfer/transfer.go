package chattransfer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	api "gf-chat/api/backend/v1"
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
	"github.com/gogf/gf/v2/util/gconv"
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
	status := "待接入"
	if relation.FormAdmin != nil {
		formName = relation.FormAdmin.Username
	}
	if relation.ToAdmin != nil {
		toName = relation.ToAdmin.Username
	}
	if relation.User != nil {
		username = relation.User.Username
	}
	if relation.CanceledAt != nil {
		status = "已取消"
	}
	if relation.AcceptedAt != nil {
		status = "已接入"
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
		Status:        status,
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
	err = s.RemoveUser(ctx, transfer.CustomerId, transfer.UserId)
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
	return s.RemoveUser(ctx, transfer.CustomerId, transfer.UserId)
}

// Create 创建转接
func (s *sChatTransfer) Create(ctx context.Context, fromAdmin *model.CustomerAdmin, toId uint, userId uint, remark string) (err error) {
	_, err = service.User().First(ctx, do.Users{
		CustomerId: fromAdmin.CustomerId,
		Id:         userId,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, "无效的用户")
	}
	_, err = service.Admin().First(ctx, do.CustomerAdmins{
		CustomerId: fromAdmin.CustomerId,
		Id:         toId,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, "无效的客服")
	}
	isValid, err := service.Chat().IsUserValid(ctx, fromAdmin.Id, userId)
	if err != nil {
		return err
	}
	if !isValid {
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, "用户已失效，无法转接")
	}
	session, err := service.ChatSession().FirstActive(ctx, userId, fromAdmin.Id, nil)
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
			UserId:     userId,
		},
	}
	newSession, err = service.ChatSession().Insert(ctx, newSession)
	if err != nil {
		return err
	}
	transfer := &model.CustomerChatTransfer{
		CustomerChatTransfers: entity.CustomerChatTransfers{
			UserId:        userId,
			FromSessionId: session.Id,
			ToSessionId:   newSession.Id,
			CustomerId:    session.CustomerId,
			FromAdminId:   fromAdmin.Id,
			ToAdminId:     toId,
			Remark:        remark,
		},
	}
	_, err = service.ChatTransfer().Save(ctx, transfer)
	if err != nil {
		return
	}
	err = s.addUser(ctx, session.CustomerId, userId, toId)
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
func (s *sChatTransfer) RemoveUser(ctx context.Context, customerId, uid uint) error {
	_, err := g.Redis().HDel(ctx, s.userKey(customerId), gconv.String(uid))
	return err
}

func (s *sChatTransfer) GetUserTransferId(ctx context.Context, customerId, uid uint) (id uint, error error) {
	val, err := g.Redis().HGet(ctx, s.userKey(customerId), gconv.String(uid))
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

// addUser 添加用户到转接列表中
func (s *sChatTransfer) addUser(ctx context.Context, customerId, uid, adminId uint) (err error) {
	_, err = g.Redis().HSet(ctx, s.userKey(customerId), g.Map{
		gconv.String(uid): adminId,
	})
	return
}

func (s *sChatTransfer) FirstActive(ctx context.Context, where any) (transfer *model.CustomerChatTransfer, err error) {
	err = s.Dao.Ctx(ctx).WhereNull("canceled_at").WhereNull("accepted_at").Where(where).Scan(&transfer)
	if err != nil {
		return
	}
	if transfer == nil {
		err = sql.ErrNoRows
		return
	}
	return
}
