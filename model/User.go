package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique_index"`
	Password string `gorm:"type:varchar(1024)"`
}
