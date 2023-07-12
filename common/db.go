package common

import (
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetDatabase(_db *gorm.DB) {
	DB = _db
}
