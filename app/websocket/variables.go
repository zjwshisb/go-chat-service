package websocket

import "ws/app/repositories"

type Where = *repositories.Where

var (
	messageRepo = &repositories.MessageRepo{}
	sessionRepo = &repositories.ChatSessionRepo{}
	autoRuleRepo = &repositories.AutoRuleRepo{}
	adminRepo = &repositories.AdminRepo{}
)
