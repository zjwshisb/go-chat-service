package repositories

import (
	"ws/app/models"
)

type userRepo struct {
	Repository[models.User]
}
