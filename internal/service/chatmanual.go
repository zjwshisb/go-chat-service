// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

type (
	IChatManual interface {
		// Add 加入到待人工接入sortSet
		Add(uid uint, gid uint) error
		// IsIn 是否在待人工接入列表中
		IsIn(uid uint, customerId uint) bool
		// Remove 从待人工接入列表中移除
		Remove(uid uint, customerId uint) error
		// GetTotalCount 获取待人工接入的数量
		GetTotalCount(customerId uint) uint
		// GetCountByTime 获取指定时间的数量
		GetCountByTime(customerId uint, min string, max string) uint
		// GetByTime 通过加入时间获取
		GetByTime(customerId uint, min string, max string) []string
		// GetTime 获取加入时间
		GetTime(uid uint, customerId uint) float64
		// GetAll 获取所有待人工接入ids
		GetAll(customerId uint) []uint
		GetBySource(customerId uint, Offset uint, count uint) []uint
	}
)

var (
	localChatManual IChatManual
)

func ChatManual() IChatManual {
	if localChatManual == nil {
		panic("implement not found for interface IChatManual, forgot register?")
	}
	return localChatManual
}

func RegisterChatManual(i IChatManual) {
	localChatManual = i
}
