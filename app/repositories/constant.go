package repositories

import "ws/app/models"

var (
	AdminRepo       = &adminRepo{}
	AutoMessageRepo = &autoMessageRepo{}
	AutoRuleRepo    = &autoRuleRepo{}
	ChatSettingRepo = &Repository[models.ChatSetting]{}
	MessageRepo     = &messageRepo{}
	ChatSessionRepo = &chatSessionRepo{}
	TransferRepo    = &transferRepo{}
	UserRepo        = &userRepo{}
)
