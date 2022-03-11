package chat

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
	"ws/app/databases"
	"ws/app/models"
	"ws/app/repositories"
)


const (
	// 转接待接入的用户 sets
	transferUserKey = "user:transfer"
)
var TransferService = &transferService{}

type transferService struct {

}

// Cancel 取消待接入的转接
func (transferService *transferService) Cancel(transfer *models.ChatTransfer) error {
	repositories.ChatSessionRepo.Delete([]*repositories.Where{
		{
			Filed: "admin_id = ?",
			Value: 0,
		},
		{
			Filed: "type = ? ",
			Value: models.ChatSessionTypeTransfer,
		},
		{
			Filed: "user_id = ?",
			Value: transfer.UserId,
		},
	})
	transfer.IsCanceled = true
	t := time.Now()
	transfer.CanceledAt = t.Unix()
	_ = repositories.TransferRepo.Save(transfer)
	_ = transferService.RemoveUser(transfer.UserId)
	return nil
}

// Create 创建转接
func (transferService *transferService) Create(fromId int64, toId int64, uid int64, remark  string) error  {
	session := repositories.ChatSessionRepo.FirstActiveByUser(uid, fromId)
	SessionService.Close(session.Id, true, true)
	if session == nil {
		return errors.New("invalid user")
	}
	now := time.Now()
	newSession := repositories.ChatSessionRepo.Create(uid, session.GroupId,models.ChatSessionTypeTransfer)
	transfer := &models.ChatTransfer{
		UserId:      uid,
		SessionId:   newSession.Id,
		FromAdminId: fromId,
		ToAdminId:   toId,
		Remark:      remark,
		CreatedAt:   now.Unix(),
	}
	repositories.TransferRepo.Save(transfer)
	_ = transferService.AddUser(uid, toId)
	return nil
}

// RemoveUser 在转接列表中移除user
func (transferService *transferService) RemoveUser(uid int64) error  {
	ctx := context.Background()
	cmd := databases.Redis.HDel(ctx, transferUserKey, strconv.FormatInt(uid, 10))
	return cmd.Err()
}

func (transferService *transferService) GetUserTransferId(uid int64) int64  {
	ctx := context.Background()
	cmd := databases.Redis.HGet(ctx, transferUserKey, strconv.FormatInt(uid, 10))
	if cmd.Err() == redis.Nil {
		return 0
	}
	adminId, _ := strconv.ParseInt(cmd.Val(), 10, 64)
	return adminId
}

// AddUser 添加用户到转接列表中
func (transferService *transferService) AddUser(uid int64, adminId int64) error {
	ctx := context.Background()
	cmd := databases.Redis.HSet(ctx, transferUserKey, uid, adminId)
	return cmd.Err()
}
