// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"github.com/gogf/gf/v2/os/gctx"
)

type (
	IChatRelation interface {
		// AddUser 接入user
		AddUser(ctx gctx.Ctx, adminId uint, uid uint) (err error)
		// UpdateUser 更新user
		// 更新limit time
		// 更新最后聊天时间
		UpdateUser(ctx gctx.Ctx, adminId uint, uid uint) (err error)
		// RemoveUser 移除user
		RemoveUser(ctx gctx.Ctx, adminId uint, uid uint) (err error)
		// IsUserValid 检查用户对于客服是否合法
		IsUserValid(ctx gctx.Ctx, adminId uint, uid uint) bool
		// IsUserExist user是否存在
		IsUserExist(ctx gctx.Ctx, adminId uint, uid uint) bool
		// GetLastChatTime 获取最后聊天时间
		GetLastChatTime(ctx gctx.Ctx, adminId uint, uid uint) uint
		// RemoveLastChatTime 移除最后聊天时间
		RemoveLastChatTime(ctx gctx.Ctx, adminId uint, uid uint) error
		// UpdateLastChatTime 更新最后聊天时间
		UpdateLastChatTime(ctx gctx.Ctx, adminId uint, uid uint) error
		// GetActiveCount 获取有效的用户数量
		GetActiveCount(ctx gctx.Ctx, adminId uint) uint
		// UpdateLimitTime 更新有效期
		UpdateLimitTime(ctx gctx.Ctx, adminId uint, uid uint, duration int64) error
		// GetLimitTime 获取有效期
		GetLimitTime(ctx gctx.Ctx, adminId uint, uid uint) int64
		GetInvalidUsers(ctx gctx.Ctx, adminId uint) []uint
		// GetUsersWithLimitTime 获取所有user以及对应的有效期
		GetUsersWithLimitTime(ctx gctx.Ctx, adminId uint) (uids []uint, times []int64)
		// SetUserAdmin SetAdmin 设置用户客服
		SetUserAdmin(ctx gctx.Ctx, uid uint, adminId uint) (err error)
		// RemoveUserAdmin RemoveAdmin 移除用户客服
		RemoveUserAdmin(ctx gctx.Ctx, uid uint) (err error)
		// GetUserValidAdmin GetValidAdmin 获取用户客服
		GetUserValidAdmin(ctx gctx.Ctx, uid uint) (uint, error)
	}
)

var (
	localChatRelation IChatRelation
)

func ChatRelation() IChatRelation {
	if localChatRelation == nil {
		panic("implement not found for interface IChatRelation, forgot register?")
	}
	return localChatRelation
}

func RegisterChatRelation(i IChatRelation) {
	localChatRelation = i
}
