package repositories

import (
	"ws/app/models"
)

type transferRepo struct {
	Repository[models.ChatTransfer]
}
