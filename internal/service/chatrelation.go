// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import "context"

type ctx = context.Context

type (
	IChatRelation interface {
		// AddUser 接入user
		AddUser(ctx ctx, adminId uint, uid uint) error
		// UpdateUser 更新user
		// 更新limit time
		// 更新最后聊天时间
		UpdateUser(ctx ctx, adminId uint, uid uint) error
		// RemoveUser 移除user
		RemoveUser(ctx ctx, adminId uint, uid uint) error
		// IsUserValid 检查用户对于客服是否合法
		IsUserValid(ctx ctx, adminId uint, uid uint) bool
		// IsUserExist user是否存在
		IsUserExist(ctx ctx, adminId uint, uid uint) bool
		// GetLastChatTime 获取最后聊天时间
		GetLastChatTime(ctx ctx, adminId uint, uid uint) uint
		// RemoveLastChatTime 移除最后聊天时间
		RemoveLastChatTime(ctx ctx, adminId uint, uid uint) error
		// UpdateLastChatTime 更新最后聊天时间
		UpdateLastChatTime(ctx ctx, adminId uint, uid uint) error
		// GetActiveCount 获取有效的用户数量
		GetActiveCount(ctx ctx, adminId uint) uint
		// UpdateLimitTime 更新有效期
		UpdateLimitTime(ctx ctx, adminId uint, uid uint, duration int64) error
		// GetLimitTime 获取有效期
		GetLimitTime(ctx ctx, adminId uint, uid uint) int64
		GetInvalidUsers(ctx ctx, adminId uint) []uint
		// GetUsersWithLimitTime 获取所有user以及对应的有效期
		GetUsersWithLimitTime(ctx ctx, adminId uint) ([]uint, []int64)
		// SetUserAdmin SetAdmin 设置用户客服
		SetUserAdmin(ctx ctx, uid uint, adminId uint) error
		// RemoveUserAdmin RemoveAdmin 移除用户客服
		RemoveUserAdmin(ctx ctx, uid uint) error
		// GetUserValidAdmin GetValidAdmin 获取用户客服
		GetUserValidAdmin(ctx ctx, uid uint) uint
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
