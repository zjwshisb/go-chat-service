package admin

import "ws/app/repositories"

type Where = *repositories.Where

var (
	messageRepo = &repositories.MessageRepo{}
	transferRepo = &repositories.TransferRepo{}
	chatSessionRepo = &repositories.ChatSessionRepo{}
	autoMessageRepo = &repositories.AutoMessageRepo{}
	autoRuleRepo = &repositories.AutoRuleRepo{}
)
