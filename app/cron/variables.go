package cron

import "ws/app/repositories"

var (
	adminRepo = &repositories.AdminRepo{}
	sessionRepo = &repositories.ChatSessionRepo{}
	messageRepo = &repositories.MessageRepo{}
)
