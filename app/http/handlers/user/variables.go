package user

import "ws/app/repositories"

type Where = *repositories.Where

var (
	messageRepo = &repositories.MessageRepo{}
)
