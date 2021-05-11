package migrate

import (
	"golang.org/x/crypto/bcrypt"
	"ws/db"
	"ws/internal/models"
)

func Seed() {
	password := []byte("user1")
	hash, _ := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	db.Db.Create(&models.User{
		Username: "user1",
		Password: string(hash),
	})
}
