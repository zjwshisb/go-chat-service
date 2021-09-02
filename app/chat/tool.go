package chat

import "ws/app/repositories"

var (
	transferRepo = &repositories.TransferRepo{}
	chatSessionRepo = &repositories.ChatSessionRepo{}
)
