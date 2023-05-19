package database

import (
	"gorm.io/gorm"
)

var DB *gorm.DB

func Set(db *gorm.DB) {
	DB = db
}
