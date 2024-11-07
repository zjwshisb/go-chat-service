package chattransfer

import (
	"context"
	"database/sql"
	"fmt"
	"gf-chat/internal/consts"
	"gf-chat/internal/dao"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/model/entity"
	"gf-chat/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
)

func init() {
	service.RegisterChatTransfer(&sChatTransfer{})
}

const (
	// 转接待接入的用户 sets
	transferUserKey = "user:%d:transfer"
)

type sChatTransfer struct {
}

func (s *sChatTransfer) Paginate(ctx context.Context, w *do.CustomerChatTransfers, p model.QueryInput) (res []*model.CustomerChatTransfer, total uint) {
	query := dao.CustomerChatTransfers.Ctx(ctx).Where(w)
	if p.WithTotal {
		i, _ := query.Clone().Count()
		total = uint(i)
		if total == 0 {
			return
		}
	}
	err := query.WithAll().Page(p.GetPage(), p.GetSize()).OrderDesc("id").Scan(&res)
	if err == sql.ErrNoRows {
		return
	}
	return
}

func (s *sChatTransfer) First(w any, with ...any) (item *model.CustomerChatTransfer, err error) {
	err = dao.CustomerChatTransfers.Ctx(gctx.New()).Where(w).With(with...).Scan(&item)
	if err != nil {
		return
	}
	if item == nil {
		err = sql.ErrNoRows
	}
	return
}

func (s *sChatTransfer) All(w any, with ...any) []*model.CustomerChatTransfer {
	res := make([]*model.CustomerChatTransfer, 0)
	dao.CustomerChatTransfers.Ctx(gctx.New()).With(with...).Where(w).Scan(&res)
	return res
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
	service.ChatSession().Close(session, true, true)
	newSession := &entity.CustomerChatSessions{
		UserId:     uid,
		QueriedAt:  gtime.New(),
		CustomerId: session.CustomerId,
		AdminId:    toId,
		Type:       consts.ChatSessionTypeTransfer,
	}
	service.ChatSession().SaveEntity(ctx, newSession)
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
