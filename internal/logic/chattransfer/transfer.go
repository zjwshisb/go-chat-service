package chattransfer

import (
	"context"
	"fmt"
	"gf-chat/internal/consts"
	"gf-chat/internal/dao"
	"gf-chat/internal/model"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"
	"gf-chat/internal/trait"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
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

func (s *sChatTransfer) ToChatTransfer(relation *model.CustomerChatTransfer) model.ChatTransfer {
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
	return model.ChatTransfer{
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
func (s *sChatTransfer) Cancel(transfer *model.CustomerChatTransfer) error {
	if transfer.ToSession != nil {
		transfer.ToSession.CanceledAt = gtime.New()
		dao.CustomerChatSessions.Ctx(gctx.New()).Save(transfer.ToSession)
	}
	transfer.CanceledAt = gtime.New()
	dao.CustomerChatTransfers.Ctx(gctx.New()).Save(transfer)
	_ = s.removeUser(transfer.CustomerId, transfer.UserId)
	service.Chat().NoticeTransfer(transfer.CustomerId, transfer.ToAdminId)
	return nil
}

func (s *sChatTransfer) Accept(transfer *model.CustomerChatTransfer) error {
	transfer.AcceptedAt = gtime.New()
	dao.CustomerChatTransfers.Ctx(gctx.New()).Save(transfer)
	service.Chat().NoticeTransfer(transfer.CustomerId, transfer.ToAdminId)
	return s.removeUser(transfer.CustomerId, transfer.UserId)
}

// Create 创建转接
func (s *sChatTransfer) Create(ctx context.Context, fromAdminId, toId, uid uint, remark string) error {
	session, err := service.ChatSession().ActiveOne(ctx, uid, fromAdminId, nil)
	if err != nil {
		return err
	}
	if session == nil {
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, "用户已失效，无法转接")
	}
	service.ChatSession().Close(ctx, session, true, false)
	newSession := &model.CustomerChatSession{
		CustomerChatSessions: entity.CustomerChatSessions{
			QueriedAt:  gtime.New(),
			CustomerId: session.CustomerId,
			AdminId:    toId,
			Type:       consts.ChatSessionTypeTransfer,
			UserId:     uid,
		},
	}

	service.ChatSession().Insert(ctx, newSession)
	transfer := &entity.CustomerChatTransfers{
		UserId:        uid,
		FromSessionId: session.Id,
		ToSessionId:   newSession.Id,
		CustomerId:    session.CustomerId,
		FromAdminId:   fromAdminId,
		ToAdminId:     toId,
		Remark:        remark,
	}
	dao.CustomerChatTransfers.Ctx(ctx).Save(transfer)
	s.addUser(session.CustomerId, uid, toId)
	service.Chat().NoticeTransfer(session.CustomerId, toId)
	return nil
}
func (s *sChatTransfer) userKey(customerId uint) string {
	return fmt.Sprintf(transferUserKey, customerId)
}

// RemoveUser 在转接列表中移除user
func (s *sChatTransfer) removeUser(customerId, uid uint) error {
	_, err := g.Redis().Do(gctx.New(), "hDel", s.userKey(customerId), uid)
	return err
}

func (s *sChatTransfer) GetUserTransferId(customerId, uid uint) uint {
	val, err := g.Redis().Do(gctx.New(), "hGet", s.userKey(customerId), uid)
	if err != nil {
		return 0
	}
	return val.Uint()
}

// AddUser 添加用户到转接列表中
func (s *sChatTransfer) addUser(customerId, uid, adminId uint) error {
	ctx := gctx.New()
	_, err := g.Redis().Do(ctx, "hSet", s.userKey(customerId), uid, adminId)
	return err
}
