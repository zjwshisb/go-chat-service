package chattransfer

import (
	"context"
	"database/sql"
	"fmt"
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
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
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

func (s *sChatTransfer) Paginate(ctx context.Context, w *do.CustomerChatTransfers, p model.QueryInput) (res []*relation.CustomerChatTransfer, total int) {
	query := dao.CustomerChatTransfers.Ctx(ctx).Where(w)
	if p.WithTotal {
		total, _ = query.Clone().Count()
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

func (s *sChatTransfer) FirstEntity(w any) *entity.CustomerChatTransfers {
	item := &entity.CustomerChatTransfers{}
	err := dao.CustomerChatTransfers.Ctx(gctx.New()).Where(w).Scan(&item)
	if err == sql.ErrNoRows {
		return nil
	}
	return item
}

func (s *sChatTransfer) FirstRelation(w any) *relation.CustomerChatTransfer {
	item := &relation.CustomerChatTransfer{}
	err := dao.CustomerChatTransfers.Ctx(gctx.New()).Where(w).WithAll().Scan(&item)
	if err == sql.ErrNoRows {
		return nil
	}
	return item
}

func (s *sChatTransfer) GetRelations(w any) []*relation.CustomerChatTransfer {
	res := make([]*relation.CustomerChatTransfer, 0)
	dao.CustomerChatTransfers.Ctx(gctx.New()).Where(w).WithAll().Scan(&res)
	return res
}

func (s *sChatTransfer) RelationToChat(relation *relation.CustomerChatTransfer) chat.Transfer {
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
	return chat.Transfer{
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
func (s *sChatTransfer) Cancel(transfer *relation.CustomerChatTransfer) error {
	if transfer.ToSession != nil {
		transfer.ToSession.CanceledAt = time.Now().Unix()
		dao.CustomerChatSessions.Ctx(gctx.New()).Save(transfer.ToSession)
	}
	transfer.CanceledAt = time.Now().Unix()
	dao.CustomerChatTransfers.Ctx(gctx.New()).Save(transfer)
	_ = s.removeUser(transfer.CustomerId, transfer.UserId)
	service.Chat().NoticeTransfer(transfer.CustomerId, transfer.ToAdminId)
	return nil
}

func (s *sChatTransfer) Accept(transfer *entity.CustomerChatTransfers) error {
	transfer.AcceptedAt = time.Now().Unix()
	dao.CustomerChatTransfers.Ctx(gctx.New()).Save(transfer)
	service.Chat().NoticeTransfer(transfer.CustomerId, transfer.ToAdminId)
	return s.removeUser(transfer.CustomerId, transfer.UserId)
}

// Create 创建转接
func (s *sChatTransfer) Create(fromAdminId, toId, uid int, remark string) error {
	session := service.ChatSession().ActiveOne(uid, fromAdminId, nil)
	if session == nil {
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, "用户已失效，无法转接")
	}
	service.ChatSession().Close(session, true, true)
	newSession := &entity.CustomerChatSessions{
		UserId:     uid,
		QueriedAt:  time.Now().Unix(),
		CustomerId: session.CustomerId,
		AdminId:    toId,
		Type:       consts.ChatSessionTypeTransfer,
	}
	service.ChatSession().SaveEntity(newSession)
	transfer := &entity.CustomerChatTransfers{
		UserId:        uid,
		FromSessionId: session.Id,
		ToSessionId:   newSession.Id,
		CustomerId:    session.CustomerId,
		FromAdminId:   gconv.Int(fromAdminId),
		ToAdminId:     toId,
		Remark:        remark,
	}
	dao.CustomerChatTransfers.Ctx(gctx.New()).Save(transfer)
	s.addUser(session.CustomerId, uid, toId)
	service.Chat().NoticeTransfer(session.CustomerId, toId)
	return nil
}
func (s *sChatTransfer) userKey(customerId int) string {
	return fmt.Sprintf(transferUserKey, customerId)
}

// RemoveUser 在转接列表中移除user
func (s *sChatTransfer) removeUser(customerId, uid int) error {
	_, err := g.Redis().Do(gctx.New(), "hDel", s.userKey(customerId), uid)
	return err
}

func (s *sChatTransfer) GetUserTransferId(customerId, uid int) int {
	val, err := g.Redis().Do(gctx.New(), "hGet", s.userKey(customerId), uid)
	if err != nil {
		return 0
	}
	return val.Int()
}

// AddUser 添加用户到转接列表中
func (s *sChatTransfer) addUser(customerId, uid, adminId int) error {
	ctx := gctx.New()
	_, err := g.Redis().Do(ctx, "hSet", s.userKey(customerId), uid, adminId)
	return err
}
